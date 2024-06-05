package database

import (
	"time"
)

// Eskimoe is a Terminal Based Chat Application and these are the models for the database.
// These models remain on the server. Eskimoe servers are like Discord servers, except
// they are self-hosted and a Eskimoe client connects to multiple Eskimoe servers.

type Permission string

const (
	ViewRoom           Permission = "view_room"
	SendMessage        Permission = "send_message"
	AddLink            Permission = "add_link"
	AddFile            Permission = "add_file"
	AddReaction        Permission = "add_reaction"
	CreatePoll         Permission = "create_poll"
	DeleteMessage      Permission = "delete_message"
	ManageRoles        Permission = "manage_roles"
	ChangeName         Permission = "change_name"
	MuteMembers        Permission = "mute_members"
	KickMembers        Permission = "kick_members"
	BanMembers         Permission = "ban_members"
	ManageRooms        Permission = "manage_rooms"
	RunCommands        Permission = "run_commands"
	ViewLogs           Permission = "view_logs"
	ViewMessageHistory Permission = "view_message_history"
	CreateEvents       Permission = "create_events"
	ManageEvents       Permission = "manage_events"
	Administrator      Permission = "administrator"
)

// Room Types: Announcement, Text, Commands, Archive
type RoomType string

const (
	Announcement RoomType = "announcement"
	Text         RoomType = "text"
	Commands     RoomType = "commands"
	Archive      RoomType = "archive"
)

type Member struct {
	ID          int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	UniqueID    string    `gorm:"unique" json:"unique_id"`
	UniqueToken string    `gorm:"unique" json:"-"`
	DisplayName string    `json:"display_name"`
	About       string    `json:"about"`
	JoinedAt    string    `json:"joined_at"`
	Pronouns    string    `json:"pronouns"`
	Roles       []Role    `gorm:"many2many:member_roles" json:"roles"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type Role struct {
	ID          int                `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name        string             `gorm:"unique" json:"name"`
	Permissions []ServerPermission `json:"permissions"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
}

// Rooms are like Channels in Discord. They are chat rooms where users can chat with each other.
type Room struct {
	ID          int              `gorm:"primaryKey;autoIncrement=true" json:"-"` // Primary Key
	Name        string           `gorm:"unique" json:"name"`
	Description string           `json:"description"`
	Type        RoomType         `json:"type"`
	Permissions []RoomPermission `json:"permissions"`
	CreatedAt   time.Time        `json:"-"`
	UpdatedAt   time.Time        `json:"-"`
}

// RoomPermissions are the permissions of a role in a room. They are used to override the server permissions.
type RoomPermission struct {
	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
	RoomID     int        `json:"-"`
	RoleID     uint       `json:"-"`
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `json:"permission"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

// Permissions of a role on the whole server in general.
type ServerPermission struct {
	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
	RoleID     uint       `json:"-"`
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `json:"permission"`
	CreatedAt  time.Time  `json:"-"`
}
