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

type ServerMode string

const (
	InviteOnly ServerMode = "invite_only"
	Open       ServerMode = "open"
	Passphrase ServerMode = "passphrase"
)

type Server struct {
	ID             int              `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name           string           `gorm:"unique" json:"name"`
	Message        string           `json:"message"`
	PublicURL      string           `json:"public_url"`
	Mode           ServerMode       `gorm:"default:open" json:"mode"`
	Members        []Member         `gorm:"many2many:server_members" json:"members"`
	Roles          []Role           `gorm:"foreignKey:ServerID" json:"roles"`
	RoomCategories []RoomCategory   `gorm:"foreignKey:ServerID" json:"room_categories"`
	Reactions      []ServerReaction `gorm:"foreignKey:ServerID" json:"reactions"`
	CategoryOrder  datatypes.JSON   `json:"category_order"`
	RoleOrder      datatypes.JSON   `json:"role_order"`
	CreatedAt      time.Time        `json:"-"`
	UpdatedAt      time.Time        `json:"-"`
}

type ServerReaction struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	ServerID  int       `json:"-"`
	Reaction  string    `json:"reaction"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type ServerInvite struct {
	ID         int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	MemberID   int       `json:"-"`
	Member     Member    `gorm:"foreignKey:MemberID" json:"member"`
	InviteCode string    `json:"invite_code"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type ServerPermission struct {
	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
	RoleID     int        `json:"-"`
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `json:"permission"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

type Member struct {
	ID                   int            `gorm:"primaryKey;autoIncrement=true" json:"-"`
	UniqueID             string         `gorm:"unique" json:"-"`
	AuthToken            string         `json:"auth_token,omitempty"`
	UniqueToken          string         `json:"-"`
	DisplayName          string         `json:"display_name"`
	About                string         `json:"about"`
	Pronouns             string         `json:"pronouns"`
	Roles                []Role         `gorm:"many2many:member_roles" json:"roles"`
	InviteCode           string         `json:"-"`
	GeneratedInviteCodes []ServerInvite `gorm:"foreignKey:MemberID" json:"-"`
	JoinedAt             time.Time      `json:"joined_at"`
	CreatedAt            time.Time      `json:"-"`
	UpdatedAt            time.Time      `json:"-"`
}

type Role struct {
	ID          int                `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name        string             `gorm:"unique" json:"name"`
	ServerID    int                `json:"-"`
	Permissions []ServerPermission `gorm:"foreignKey:RoleID" json:"permissions"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
}

type RoomCategory struct {
	ID          int            `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name        string         `gorm:"unique" json:"name"`
	Description string         `json:"description"`
	ServerID    int            `json:"-"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	Rooms       []Room         `gorm:"foreignKey:CategoryID" json:"rooms"`
	RoomOrder   datatypes.JSON `json:"room_order"`
}

type Room struct {
	ID          int              `gorm:"primaryKey;autoIncrement=true" json:"-"`
	Name        string           `gorm:"unique" json:"name"`
	Description string           `json:"description"`
	Type        RoomType         `json:"type"`
	CategoryID  int              `json:"-"`
	Permissions []RoomPermission `gorm:"foreignKey:RoomID" json:"permissions"`
	Messages    []Message        `gorm:"foreignKey:RoomID" json:"messages"`
	CreatedAt   time.Time        `json:"-"`
	UpdatedAt   time.Time        `json:"-"`
}

type RoomPermission struct {
	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
	RoomID     int        `json:"-"`
	RoleID     int        `json:"-"`
	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
	Permission Permission `json:"permission"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
}

type Message struct {
	ID          int                 `gorm:"primaryKey;autoIncrement=true" json:"-"`
	AuthorID    int                 `json:"-"`
	Author      Member              `gorm:"foreignKey:AuthorID" json:"author"`
	RoomID      int                 `json:"-"`
	Content     string              `json:"content"`
	Reactions   []MessageReaction   `gorm:"foreignKey:MessageID" json:"reactions"`
	Attachments []MessageAttachment `gorm:"foreignKey:MessageID" json:"attachments"`
	Edited      bool                `json:"edited"`
	CreatedAt   time.Time           `json:"-"`
	UpdatedAt   time.Time           `json:"-"`
}

type MessageAttachment struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	MessageID int       `json:"-"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MessageReaction struct {
	ID        int       `gorm:"primaryKey;autoIncrement=true" json:"-"`
	MessageID int       `json:"-"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message"`
	Reaction  string    `json:"reaction"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// type Server struct {
// 	ID             int               `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	Name           string            `gorm:"unique" json:"name"`
// 	Message        string            `json:"message"`
// 	PublicURL      string            `json:"public_url"`
// 	Mode           ServerMode        `json:"mode" gorm:"default:open"`
// 	Members        []Member          `gorm:"many2many:server_members" json:"members"`
// 	Roles          []Role            `gorm:"foreignKey:ServerID" json:"roles"`
// 	RoomCategories []RoomCategory    `gorm:"foreignKey:ServerID" json:"room_categories"`
// 	Reactions      []ServerReactions `gorm:"foreignKey:ServerID" json:"reactions"`
// 	CategoryOrder  []int             `json:"category_order"`
// 	RoleOrder      []int             `json:"role_order"`
// 	CreatedAt      time.Time         `json:"-"`
// 	UpdatedAt      time.Time         `json:"-"`
// }

// type ServerReactions struct {
// 	ID        int    `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	ServerID  int    `json:"-"`
// 	Reaction  string `json:"reaction"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

// type ServerInvites struct {
// 	ID         int    `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	MemberID   int    `json:"-"`
// 	Member     Member `gorm:"foreignKey:MemberID" json:"member"`
// 	InviteCode string `json:"invite_code"`
// 	CreatedAt  time.Time
// 	UpdatedAt  time.Time
// }

// type ServerPermissions struct {
// 	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	RoleID     uint       `json:"-"`
// 	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
// 	Permission Permission `json:"permission"`
// 	CreatedAt  time.Time  `json:"-"`
// }

// type Member struct {
// 	ID                   int             `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	UniqueID             string          `gorm:"unique_id" json:"-"`
// 	AuthToken            string          `gorm:"auth_token" json:"auth_token,omitempty"`
// 	UniqueToken          string          `gorm:"unique_token" json:"-"`
// 	DisplayName          string          `json:"display_name"`
// 	About                string          `json:"about"`
// 	Pronouns             string          `json:"pronouns"`
// 	Roles                []Role          `gorm:"many2many:member_roles" json:"roles"`
// 	InviteCode           string          `json:"-"`
// 	GeneratedInviteCodes []ServerInvites `gorm:"foreignKey:MemberID" json:"-"`
// 	JoinedAt             time.Time       `json:"joined_at"`
// 	CreatedAt            time.Time       `json:"-"`
// 	UpdatedAt            time.Time       `json:"-"`
// }

// type Role struct {
// 	ID          int                 `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	Name        string              `gorm:"unique" json:"name"`
// 	Permissions []ServerPermissions `json:"permissions"`
// 	CreatedAt   time.Time           `json:"-"`
// 	UpdatedAt   time.Time           `json:"-"`
// }

// type RoomCategory struct {
// 	ID          int    `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	Name        string `gorm:"unique" json:"name"`
// 	Description string `json:"description"`
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// 	Rooms       []Room `gorm:"foreignKey:CategoryID" json:"rooms"`
// 	RoomOrder   []int  `json:"room_order"`
// }

// type Room struct {
// 	ID          int               `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	Name        string            `gorm:"unique" json:"name"`
// 	Description string            `json:"description"`
// 	Type        RoomType          `json:"type"`
// 	Permissions []RoomPermissions `json:"permissions"`
// 	Messages    []Message         `json:"messages"`
// 	CreatedAt   time.Time         `json:"-"`
// 	UpdatedAt   time.Time         `json:"-"`
// }

// type RoomPermissions struct {
// 	ID         int        `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	RoomID     int        `json:"-"`
// 	RoleID     uint       `json:"-"`
// 	Role       Role       `gorm:"foreignKey:RoleID" json:"role"`
// 	Permission Permission `json:"permission"`
// 	CreatedAt  time.Time  `json:"-"`
// 	UpdatedAt  time.Time  `json:"-"`
// }

// type Message struct {
// 	ID          int                 `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	Author      Member              `gorm:"foreignKey:AuthorID" json:"author"`
// 	AuthorID    int                 `json:"-"`
// 	Content     string              `json:"content"`
// 	Reactions   []MessageReaction   `json:"reactions"`
// 	Attachments []MessageAttachment `json:"attachments"`
// 	Edited      bool                `json:"edited"`
// 	CreatedAt   time.Time           `json:"-"`
// 	UpdatedAt   time.Time           `json:"-"`
// }

// type MessageAttachment struct {
// 	ID        int     `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	MessageID int     `json:"-"`
// 	Message   Message `gorm:"foreignKey:MessageID" json:"message"`
// 	URL       string  `json:"url"`
// 	Type      string  `json:"type"`
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

// type MessageReaction struct {
// 	ID        int             `gorm:"primaryKey;autoIncrement=true" json:"-"`
// 	MessageID int             `json:"-"`
// 	Message   Message         `gorm:"foreignKey:MessageID" json:"message"`
// 	AuthorID  int             `json:"-"`
// 	Author    Member          `gorm:"foreignKey:AuthorID" json:"author"`
// 	Reaction  ServerReactions `json:"reaction"`
// 	Reactors  []Member        `gorm:"many2many:reactors" json:"reactors"`
// 	CreatedAt time.Time       `json:"-"`
// 	UpdatedAt time.Time       `json:"-"`
// }
