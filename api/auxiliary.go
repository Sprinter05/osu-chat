package api

import "time"

// Time according to ISO 8601

func ParseTime(s string) (time.Time, error) {
	return time.Parse(TIME_FORMAT, s)
}

func FormatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

// Channel types

func ChannelTypeToString(c ChannelType) string {
	switch c {
	case PUBLIC:
		return "PUBLIC"
	case PRIVATE:
		return "PRIVATE"
	case MULTIPLAYER:
		return "MULTIPLAYER"
	case SPECTATOR:
		return "SPECTATOR"
	case TEMPORARY:
		return "TEMPORARY"
	case PM:
		return "PM"
	case GROUP:
		return "GROUP"
	case ANNOUNCE:
		return "ANNOUNCE"
	}

	return "UNKNOWN"
}

func StringToChannelType(s string) ChannelType {
	switch s {
	case "PUBLIC":
		return PUBLIC
	case "PRIVATE":
		return PRIVATE
	case "MULTIPLAYER":
		return MULTIPLAYER
	case "SPECTATOR":
		return SPECTATOR
	case "TEMPORARY":
		return TEMPORARY
	case "PM":
		return PM
	case "GROUP":
		return GROUP
	case "ANNOUNCE":
		return ANNOUNCE
	}

	return UNKNOWN
}
