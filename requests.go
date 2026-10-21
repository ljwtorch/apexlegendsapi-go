package apex

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Platform identifies an upstream platform. Allowed values depend on the method.
type Platform string

const (
	PlatformPC     Platform = "PC"
	PlatformPS4    Platform = "PS4"
	PlatformX1     Platform = "X1"
	PlatformSwitch Platform = "SWITCH"
	PlatformAny    Platform = "ANY"
)

// GameMode filters match history.
type GameMode string

const (
	ModeBattleRoyale GameMode = "BATTLE_ROYALE"
	ModeArenas       GameMode = "ARENAS"
	ModeUnknown      GameMode = "UNKNOWN"
)

// HistoryAction selects an operation on the legacy match tracking list.
type HistoryAction string

const (
	HistoryInfo   HistoryAction = "info"
	HistoryGet    HistoryAction = "get"
	HistoryAdd    HistoryAction = "add"
	HistoryDelete HistoryAction = "delete"
)

// PlayerOptions controls optional player statistics parameters.
type PlayerOptions struct {
	// Version 0 omits the parameter (upstream default: 5). Valid: 1, 2, 4, 5.
	Version int
	// EnableClubsBeta nil omits the parameter; a pointer sends true or false.
	EnableClubsBeta *bool
	// Presence-based switches are sent only when true.
	SkipRank     bool
	Merge        bool
	RemoveMerged bool
}

// PlayerByNameRequest queries a player by name. PC requires the Origin name,
// including for Steam players. Switch is supported only by the UID method.
type PlayerByNameRequest struct {
	Player   string
	Platform Platform
	Options  PlayerOptions
}

// PlayerByUIDRequest queries a stable player UID; Switch is supported here.
type PlayerByUIDRequest struct {
	UID      string
	Platform Platform
	Options  PlayerOptions
}

// GamesRequest queries collected matches. Nil numeric options are omitted;
// non-nil pointers send their values, including zero. Timestamps are Unix epoch
// integers, as accepted by the upstream API.
type GamesRequest struct {
	UID   string
	Mode  GameMode
	Start *int64
	End   *int64
	Limit *int
}

// LegacyHistoryRequest operates on the legacy tracking list. UID and Platform
// are required for all actions, following the published parameter table.
type LegacyHistoryRequest struct {
	UID      string
	Platform Platform
	Action   HistoryAction
}

// LeaderboardRequest queries a leaderboard. Legend starts with a capital letter;
// use Global for global trackers. Empty Key requests the available tracker keys.
type LeaderboardRequest struct {
	Platform Platform
	Legend   string
	Key      string
}

// MapRotationOptions selects version 1 (BR pubs) or 2 (all modes).
// Version 0, or nil options, preserves the upstream default (version 1).
type MapRotationOptions struct{ Version int }

// NameToUIDRequest resolves a name on PC (Origin), PlayStation, or Xbox.
type NameToUIDRequest struct {
	Player   string
	Platform Platform
}

func required(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("apex: %s is required", name)
	}
	return nil
}

func platformQuery(name, value string, platform Platform, allowSwitch, allowAny bool) (url.Values, error) {
	if err := required(name, value); err != nil {
		return nil, err
	}
	switch platform {
	case PlatformPC, PlatformPS4, PlatformX1:
	case PlatformSwitch:
		if !allowSwitch {
			return nil, errors.New("apex: SWITCH is not supported by this method")
		}
	case PlatformAny:
		if !allowAny {
			return nil, errors.New("apex: ANY is not supported by this method")
		}
	default:
		return nil, errors.New("apex: invalid platform")
	}
	return url.Values{name: {value}, "platform": {string(platform)}}, nil
}

func (o PlayerOptions) apply(q url.Values) error {
	switch o.Version {
	case 0:
	case 1, 2, 4, 5:
		q.Set("version", strconv.Itoa(o.Version))
	default:
		return errors.New("apex: player version must be 1, 2, 4, or 5")
	}
	if o.EnableClubsBeta != nil {
		q.Set("enableClubsBeta", strconv.FormatBool(*o.EnableClubsBeta))
	}
	if o.SkipRank {
		q.Set("skipRank", "1")
	}
	if o.Merge {
		q.Set("merge", "1")
	}
	if o.RemoveMerged {
		q.Set("removeMerged", "1")
	}
	return nil
}
