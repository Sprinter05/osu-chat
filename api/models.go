package api

import "time"

/* IDENTIFICATORS */

type UserId int
type ChannelId int
type MessageId int

/* ENUMERATIONS */

type ChannelType string

const (
	UNKNOWN     ChannelType = "UNKNOWN"
	PUBLIC      ChannelType = "PUBLIC"
	PRIVATE     ChannelType = "PRIVATE"
	MULTIPLAYER ChannelType = "MULTIPLAYER"
	SPECTATOR   ChannelType = "SPECTATOR"
	TEMPORARY   ChannelType = "TEMPORARY"
	PM          ChannelType = "PM"
	GROUP       ChannelType = "GROUP"
	ANNOUNCE    ChannelType = "ANNOUNCE"
)

type MessageType string

const (
	ACTION    = "action"
	MARKWDOWN = "markdown"
	PLAIN     = "plain"
)

/* OBJECTS */

type User struct {
	// Constant fields
	AvatarUrl     string    `json:"avatar_url"`
	CountryCode   string    `json:"country_code"`
	DefaultGroup  string    `json:"default_group"`
	Id            UserId    `json:"id"`
	IsActive      bool      `json:"is_active"`
	IsBot         bool      `json:"is_bot"`
	IsDeleted     bool      `json:"is_deleted"`
	IsOnline      bool      `json:"is_online"`
	IsSupporter   bool      `json:"is_supporter"`
	LastVisit     time.Time `json:"last_visit"`
	PMFriendsOnly bool      `json:"pm_friends_only"`
	ProfileColour bool      `json:"profile_colour"`
	Username      string    `json:"username"`

	// Optional fields (only necessary ones)
	IsRestricted    *bool `json:"is_restricted,omitempty"`
	SessionVerified *bool `json:"session_verified,omitempty"`
	UnreadPMCount   *int  `json:"unread_pm_count,omitempty"`
}

type ChatChannelUserAttributes struct {
	CanMessage      bool      `json:"can_message"`
	CanMessageError string    `json:"can_message_error"`
	LastReadId      MessageId `json:"last_read_id"`
}

type ChatChannel struct {
	// Constant fields
	ChannelId          ChannelId   `json:"channel_id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Icon               string      `json:"icon"`
	Type               ChannelType `json:"type"`
	MessageLengthLimit int         `json:"message_length_limit"`
	Moderated          bool        `json:"moderated"`
	UUID               string      `json:"uuid"`

	// Optional fields
	CurrentUserAttributes *ChatChannelUserAttributes `json:"current_user_attributes,omitempty"`
	LastMessageId         *MessageId                 `json:"last_message_id"`
	Users                 []int                      `json:"users"`
}

type ChatMessage struct {
	// Constant fields
	ChannelId ChannelId   `json:"channel_id"`
	Content   string      `json:"content"`
	IsAction  bool        `json:"is_action"`
	MessageId MessageId   `json:"message_id"`
	SenderId  UserId      `json:"sender_id"`
	Timestamp time.Time   `json:"timestamp"`
	Type      MessageType `json:"type"`
	UUID      string      `json:"uuid"`

	// Optional fields
	Sender *User `json:"sender,omitempty"`
}
