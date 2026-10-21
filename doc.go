// Package apex provides a client for the unofficial Apex Legends Status API.
//
// Construct a Client with NewClient and reuse it across goroutines. All calls
// accept a context and return the complete upstream JSON in a Response. Decode
// that response into application-defined types with Response.Decode.
//
// Rate limiting, retries, match collection scheduling, and webhook handling are
// the responsibility of the application. See https://apexlegendsapi.com/ for
// access restrictions and attribution requirements.
package apex
