package apex

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the Apex Legends Status API address.
const DefaultBaseURL = "https://api.apexlegendsstatus.com"

// AuthMode selects how the API key is sent.
type AuthMode int

const (
	// AuthHeader sends Authorization: API_KEY (without a Bearer prefix).
	AuthHeader AuthMode = iota
	// AuthQuery sends the API key in the auth query parameter.
	AuthQuery
)

// Client is safe for concurrent use. Its settings are fixed at construction.
// A custom HTTP transport must also support concurrent use.
type Client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
	authMode   AuthMode
}

// Option configures a Client at construction time.
type Option func(*Client) error

// WithBaseURL overrides the API address, for example for a proxy or local tests.
// A path prefix is allowed; credentials, query strings, and fragments are not.
func WithBaseURL(address string) Option {
	return func(c *Client) error {
		u, err := url.Parse(address)
		if err != nil || u == nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" {
			return errors.New("apex: invalid base URL: expected an absolute HTTP(S) URL without credentials, query, or fragment")
		}
		c.baseURL = u
		return nil
	}
}

// WithHTTPClient supplies the HTTP client. Its timeout and redirect policy are
// used as provided. Configure it before NewClient and do not mutate it during use.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) error {
		if client == nil {
			return errors.New("apex: HTTP client must not be nil")
		}
		c.httpClient = client
		return nil
	}
}

// WithAuthMode selects header or query authentication. Header is the default.
func WithAuthMode(mode AuthMode) Option {
	return func(c *Client) error {
		if mode != AuthHeader && mode != AuthQuery {
			return errors.New("apex: invalid authentication mode")
		}
		c.authMode = mode
		return nil
	}
}

// NewClient constructs a client with a 30-second timeout and header authentication.
// Options are applied in order. No network request is made during construction.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" || strings.ContainsAny(apiKey, "\r\n") {
		return nil, errors.New("apex: API key must not be empty or contain newlines")
	}
	u, _ := url.Parse(DefaultBaseURL)
	c := &Client{apiKey: apiKey, baseURL: u, httpClient: &http.Client{Timeout: 30 * time.Second}}
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("apex: option must not be nil")
		}
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Response contains the upstream body unchanged, including unknown JSON fields.
// On HTTP errors it may contain a non-JSON body; Decode then returns an error.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       json.RawMessage
}

// Decode decodes exactly one JSON value into v. When v contains interface values,
// numbers are json.Number rather than float64 to avoid losing integer precision.
func (r *Response) Decode(v any) error {
	if r == nil {
		return errors.New("apex: cannot decode a nil response")
	}
	d := json.NewDecoder(bytes.NewReader(r.Body))
	d.UseNumber()
	if err := d.Decode(v); err != nil {
		return fmt.Errorf("apex: decode response: %w", err)
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		if err == nil {
			return errors.New("apex: response contains multiple JSON values")
		}
		return fmt.Errorf("apex: decode response: %w", err)
	}
	return nil
}

// ErrAccessDenied indicates that the upstream API rejected the API key or its
// access permissions. In particular, /games can report this with HTTP 404.
var ErrAccessDenied = errors.New("apex: API key does not exist or lacks access to this data; check API key and endpoint whitelist permissions")

// ErrPlayerNotFound indicates an upstream player lookup failure. This can also
// occur when the upstream Origin lookup service fails for an existing player.
var ErrPlayerNotFound = errors.New("apex: player lookup failed; check the Origin account name or try again later")

// APIError is returned for every non-2xx HTTP status, including undocumented ones.
// It is also returned for a JSON Error response, even with a 2xx HTTP status.
// Body and Header preserve the upstream response. Error omits the body and URL so
// it cannot accidentally print an API key reflected in either of those fields.
type APIError struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (e *APIError) Error() string {
	message := fmt.Sprintf("apex: API returned HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
	if e.StatusCode >= 200 && e.StatusCode < 300 {
		message = fmt.Sprintf("apex: API reported an error with HTTP %d %s", e.StatusCode, http.StatusText(e.StatusCode))
	}
	if cause := e.Unwrap(); cause != nil {
		message += ": " + strings.TrimPrefix(cause.Error(), "apex: ")
	}
	return message
}

// Unwrap classifies known upstream error responses independently of
// its HTTP status. Unknown bodies are never included in the printable error.
func (e *APIError) Unwrap() error {
	if strings.TrimSpace(string(e.Body)) == "Error: API key doesn't exist !" {
		return ErrAccessDenied
	}
	switch apiErrorMessage(e.Body) {
	case "This API key doesn't exist or you are not allowed to access this data.":
		return ErrAccessDenied
	case "Player not found", "Player not found.", "Player not found. Please try again (err origin lookup)":
		return ErrPlayerNotFound
	}
	return nil
}

func apiErrorMessage(raw []byte) string {
	var body struct {
		Error string `json:"Error"`
	}
	if json.Unmarshal(raw, &body) != nil {
		return ""
	}
	return strings.TrimSpace(body.Error)
}

func (c *Client) get(ctx context.Context, path string, query url.Values) (*Response, error) {
	if ctx == nil {
		return nil, errors.New("apex: context must not be nil")
	}
	u := *c.baseURL
	u.Path = strings.TrimRight(u.Path, "/") + path
	u.RawPath = ""
	if query == nil {
		query = make(url.Values)
	}
	if c.authMode == AuthQuery {
		query.Set("auth", c.apiKey)
	}
	u.RawQuery = query.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("apex: create request: %w", err)
	}
	// Extensionless upstream routes negotiate against PHP script media types
	// before executing them. A JSON-only Accept can cause HTTP 406 even though
	// the script itself returns JSON. Match the documented curl examples.
	req.Header.Set("Accept", "*/*")
	if c.authMode == AuthHeader {
		req.Header.Set("Authorization", c.apiKey)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// net/http wraps transport/context errors with the request URL. Remove it
		// from the printable wrapper while preserving the underlying error chain.
		var ue *url.Error
		if errors.As(err, &ue) {
			err = &url.Error{Op: ue.Op, URL: "[redacted]", Err: ue.Err}
		}
		return nil, fmt.Errorf("apex: request failed: %w", err)
	}
	defer resp.Body.Close()
	body, readErr := io.ReadAll(resp.Body)
	r := &Response{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: json.RawMessage(body)}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || apiErrorMessage(body) != "" {
		apiErr := &APIError{StatusCode: resp.StatusCode, Header: resp.Header.Clone(), Body: bytes.Clone(body)}
		if readErr != nil {
			return r, errors.Join(apiErr, fmt.Errorf("apex: read response: %w", readErr))
		}
		return r, apiErr
	}
	if readErr != nil {
		return r, fmt.Errorf("apex: read response: %w", readErr)
	}
	return r, nil
}
