package apex

import (
	"context"
	"errors"
	"net/url"
	"strconv"
)

// PlayerByName returns player statistics from /bridge using an account name.
func (c *Client) PlayerByName(ctx context.Context, request PlayerByNameRequest) (*Response, error) {
	q, err := platformQuery("player", request.Player, request.Platform, false, false)
	if err != nil {
		return nil, err
	}
	if err := request.Options.apply(q); err != nil {
		return nil, err
	}
	return c.get(ctx, "/bridge", q)
}

// PlayerByUID returns player statistics from /bridge using a stable UID.
func (c *Client) PlayerByUID(ctx context.Context, request PlayerByUIDRequest) (*Response, error) {
	q, err := platformQuery("uid", request.UID, request.Platform, true, false)
	if err != nil {
		return nil, err
	}
	if err := request.Options.apply(q); err != nil {
		return nil, err
	}
	return c.get(ctx, "/bridge", q)
}

// Games returns collected match history. Upstream access requires a whitelist
// and is currently unavailable to new users. This does not start collection.
func (c *Client) Games(ctx context.Context, request GamesRequest) (*Response, error) {
	if err := required("uid", request.UID); err != nil {
		return nil, err
	}
	q := url.Values{"uid": {request.UID}}
	switch request.Mode {
	case "":
	case ModeBattleRoyale, ModeArenas, ModeUnknown:
		q.Set("mode", string(request.Mode))
	default:
		return nil, errors.New("apex: invalid game mode")
	}
	if request.Start != nil {
		q.Set("start", strconv.FormatInt(*request.Start, 10))
	}
	if request.End != nil {
		q.Set("end", strconv.FormatInt(*request.End, 10))
	}
	if request.Limit != nil {
		q.Set("limit", strconv.Itoa(*request.Limit))
	}
	return c.get(ctx, "/games", q)
}

// LegacyHistory performs info/get/add/delete using /bridge?history=1.
// Add and delete modify the upstream tracking list. This API is currently
// unavailable to new users.
func (c *Client) LegacyHistory(ctx context.Context, request LegacyHistoryRequest) (*Response, error) {
	q, err := platformQuery("uid", request.UID, request.Platform, true, false)
	if err != nil {
		return nil, err
	}
	switch request.Action {
	case HistoryInfo, HistoryGet, HistoryAdd, HistoryDelete:
	default:
		return nil, errors.New("apex: invalid history action")
	}
	q.Set("history", "1")
	q.Set("action", string(request.Action))
	return c.get(ctx, "/bridge", q)
}

// Leaderboard returns leaderboard data, or available tracker keys if Key is
// empty. This API is currently unavailable to new users. Displayed data requires
// attribution to Apex Legends Status unless white-label access was granted.
func (c *Client) Leaderboard(ctx context.Context, request LeaderboardRequest) (*Response, error) {
	q, err := platformQuery("legend", request.Legend, request.Platform, true, true)
	if err != nil {
		return nil, err
	}
	if request.Key != "" {
		q.Set("key", request.Key)
	}
	return c.get(ctx, "/leaderboard", q)
}

// MapRotation returns current and next maps. Use version 2 for all modes.
func (c *Client) MapRotation(ctx context.Context, options *MapRotationOptions) (*Response, error) {
	q := make(url.Values)
	if options != nil {
		switch options.Version {
		case 0:
		case 1, 2:
			q.Set("version", strconv.Itoa(options.Version))
		default:
			return nil, errors.New("apex: map rotation version must be 1 or 2")
		}
	}
	return c.get(ctx, "/maprotation", q)
}

// Predator returns Predator thresholds and Masters counts by platform.
func (c *Client) Predator(ctx context.Context) (*Response, error) {
	return c.get(ctx, "/predator", nil)
}

// Servers returns server status. Displayed data requires an Apex Legends Status
// link or the attribution "Data from apexlegendsstatus.com".
func (c *Client) Servers(ctx context.Context) (*Response, error) {
	return c.get(ctx, "/servers", nil)
}

// OriginUID resolves an Origin name without fetching player statistics. The
// player must have played Apex Legends previously.
func (c *Client) OriginUID(ctx context.Context, player string) (*Response, error) {
	if err := required("player", player); err != nil {
		return nil, err
	}
	return c.get(ctx, "/origin", url.Values{"player": {player}})
}

// NameToUID resolves a name on PC (Origin), PlayStation, or Xbox.
func (c *Client) NameToUID(ctx context.Context, request NameToUIDRequest) (*Response, error) {
	q, err := platformQuery("player", request.Player, request.Platform, false, false)
	if err != nil {
		return nil, err
	}
	return c.get(ctx, "/nametouid", q)
}
