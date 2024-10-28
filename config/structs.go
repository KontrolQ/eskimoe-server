package config

// Socket Broadcast Struct

type BroadcastType int

const (
	MessageCreated BroadcastType = iota
	MessageDeleted
	MessageEdited
	MessageBulkDeleted
	MessageReactionCreated
	MessageReactionDeleted
	MessageReactionUpdated
	RoomCreated
	RoomDeleted
	RoomUpdated
	CategoryCreated
	CategoryDeleted
	CategoryUpdated
	CategoryOrderUpdated
	MemberJoined
	MemberLeft
	MemberBanned
	MemberKicked
	MemberUnbanned
	MemberUpdated
	RoleCreated
	RoleDeleted
	RoleUpdated
)

type SocketBroadcast struct {
	BroadcastType BroadcastType `json:"broadcast_type"`
	Data          interface{}   `json:"data"`
}
