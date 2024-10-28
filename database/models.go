package database

import (
	"time"

	"gorm.io/datatypes"
)

// Eskimoe is a Terminal Based Chat Application and these are the models for the database.
// These models remain on the server. Eskimoe servers are like Discord servers, except
// they are self-hosted and a Eskimoe client connects to multiple Eskimoe servers.

type Permission string

const (
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
	GenerateInvites    Permission = "generate_invites"
	Administrator      Permission = "administrator"
)

type MemberStatus string

const (
	Online  MemberStatus = "online"
	Idle    MemberStatus = "idle"
	Offline MemberStatus = "offline"
	Left    MemberStatus = "left"
)

// Room Types: Announcement, Text, Commands, Archive
type RoomType string

const (
	Announcement RoomType = "announcement"
	Text         RoomType = "text"
	Commands     RoomType = "commands"
	Archive      RoomType = "archive"
)

type ServerMode string

const (
	InviteOnly ServerMode = "invite_only"
	Open       ServerMode = "open"
	Passphrase ServerMode = "passphrase"
)

type LogType string

const (
	CategoryCreated    LogType = "category_created"
	CategoryDeleted    LogType = "category_deleted"
	CategoryUpdated    LogType = "category_updated"
	RoomCreated        LogType = "room_created"
	RoomDeleted        LogType = "room_deleted"
	RoomUpdated        LogType = "room_updated"
	MessageDeleted     LogType = "message_deleted"
	MessageBulkDeleted LogType = "message_bulk_deleted"
	MemberBanned       LogType = "member_banned"
	MemberKicked       LogType = "member_kicked"
	MemberUnbanned     LogType = "member_unbanned"
	MemberUpdated      LogType = "member_updated"
	RoleCreated        LogType = "role_created"
	RoleDeleted        LogType = "role_deleted"
	RoleUpdated        LogType = "role_updated"
	InviteGenerated    LogType = "invite_generated"
	InviteUsed         LogType = "invite_used"
	InviteDeleted      LogType = "invite_deleted" // Only unused invites can be deleted.
	ReactionCreated    LogType = "reaction_created"
	ReactionDeleted    LogType = "reaction_deleted"
	ReactionUpdated    LogType = "reaction_updated"
	EventCreated       LogType = "event_created"
	EventDeleted       LogType = "event_deleted"
	EventUpdated       LogType = "event_updated"
)

type Server struct {
	ID              int                      `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name            string                   `gorm:"not null" json:"name"`
	Message         string                   `json:"message"`
	PublicURL       string                   `json:"public_url"`
	Mode            ServerMode               `gorm:"not null;default:'open'" json:"mode"`
	Passphrase      string                   `json:"-"`
	Categories      []Category               `gorm:"foreignKey:ServerID" json:"categories"`
	CategoryOrder   datatypes.JSONSlice[int] `gorm:"type:json" json:"category_order"`
	ServerReactions []ServerReaction         `gorm:"foreignKey:ServerID" json:"server_reactions"`
	Invites         []Invite                 `gorm:"foreignKey:ServerID" json:"-"`
	Roles           []Role                   `gorm:"foreignKey:ServerID" json:"roles"`
	RoleOrder       datatypes.JSONSlice[int] `gorm:"type:json" json:"role_order"`
	Events          []Event                  `gorm:"foreignKey:ServerID" json:"events"`
	Logs            []Log                    `gorm:"foreignKey:ServerID" json:"logs,omitempty"`
	Members         []Member                 `gorm:"foreignKey:ServerID" json:"members"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"-"`
}

type Category struct {
	ID        int                      `gorm:"primaryKey;autoIncrement=true" json:"id"`
	Name      string                   `gorm:"not null" json:"name"`
	Rooms     []Room                   `gorm:"foreignKey:CategoryID" json:"rooms"`
	RoomOrder datatypes.JSONSlice[int] `gorm:"type:json" json:"room_order"`
	ServerID  int                      `json:"-"`
	Server    Server                   `json:"-"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"-"`
}

type Room struct {
	ID          int       `gorm:"primaryKey;autoIncrement=true" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Type        RoomType  `gorm:"not null;default:'text'" json:"type"`
	Messages    []Message `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
	CategoryID  int       `json:"-"`
	Category    Category  `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}

type Message struct {
	ID          int                 `gorm:"primaryKey;autoIncrement=true" json:"id"`
	Content     string              `gorm:"not null" json:"content"`
	AuthorID    int                 `json:"-"`
	Author      Member              `json:"author"`
	Reactions   []MessageReaction   `gorm:"foreignKey:MessageID" json:"reactions"`
	Attachments []MessageAttachment `gorm:"foreignKey:MessageID" json:"attachments"`
	Edited      bool                `json:"edited"`
	RoomID      int                 `json:"room_id"`
	Room        Room                `json:"-"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"-"`
}

type MessageReaction struct {
	ID         int            `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Reaction   ServerReaction `gorm:"foreignKey:ReactionID" json:"reaction"`
	Members    []Member       `gorm:"many2many:message_reaction_members" json:"members"`
	Count      int            `json:"count"`
	MessageID  int            `json:"-"`
	Message    Message        `json:"-"`
	ReactionID int            `json:"-"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
}

type MessageAttachment struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Type      string    `gorm:"not null" json:"type"`
	URL       string    `gorm:"not null" json:"url"`
	MessageID int       `json:"-"`
	Message   Message   `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type ServerReaction struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Reaction  string    `gorm:"not null" json:"reaction"`
	Color     string    `gorm:"not null" json:"color"`
	ServerID  int       `json:"-"`
	Server    Server    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Invite struct {
	ID            int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Code          string    `gorm:"not null" json:"code"`
	Used          bool      `json:"used"`
	GeneratedBy   Member    `gorm:"foreignKey:GeneratedByID" json:"generated_by"`
	GeneratedByID int       `json:"-"`
	UsedBy        Member    `gorm:"foreignKey:UsedByID" json:"used_by"`
	UsedByID      int       `json:"-"`
	ServerID      int       `json:"-"`
	Server        Server    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}

type Role struct {
	ID          int                             `gorm:"primaryKey;autoIncrement=true" json:"id"`
	Name        string                          `gorm:"not null" json:"name"`
	Permissions datatypes.JSONSlice[Permission] `gorm:"type:json" json:"permissions"`
	SystemRole  bool                            `json:"system_role"`
	ServerID    int                             `json:"-"`
	Server      Server                          `json:"-"`
	CreatedAt   time.Time                       `json:"created_at"`
	UpdatedAt   time.Time                       `json:"-"`
}

type Event struct {
	ID          int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CreatedBy   Member    `gorm:"foreignKey:CreatedByID" json:"created_by"`
	CreatedByID int       `json:"-"`
	Intereested []Member  `gorm:"many2many:event_interested" json:"interested"`
	ServerID    int       `json:"-"`
	Server      Server    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
}

type Log struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Type      LogType   `gorm:"not null" json:"type"`
	Content   string    `json:"content"`
	MemberID  int       `json:"-"`
	Member    Member    `json:"member"`
	ServerID  int       `json:"-"`
	Server    Server    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}

type Member struct {
	ID          int          `gorm:"primaryKey;autoIncrement=true" json:"-"`
	UniqueID    string       `gorm:"not null;unique" json:"uid"`
	AuthToken   string       `gorm:"not null" json:"-"`
	UniqueToken string       `gorm:"not null;unique" json:"-"`
	DisplayName string       `gorm:"not null" json:"display_name"`
	About       string       `json:"about"`
	Pronouns    string       `json:"pronouns"`
	InviteCode  string       `json:"-"` // This is the code that the member used to join the server.
	Roles       []Role       `gorm:"many2many:member_roles" json:"roles,omitempty"`
	ServerID    int          `json:"-"`
	Server      Server       `json:"-"`
	Status      MemberStatus `gorm:"not null;default:'online'" json:"status"`
	JoinedAt    time.Time    `json:"joined_at"`
	CreatedAt   time.Time    `json:"-"`
	UpdatedAt   time.Time    `json:"-"`
}
