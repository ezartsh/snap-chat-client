package models

import "github.com/google/uuid"

const (
	MT_Login       = "mt_login"        // message type for act login
	MT_Register    = "mt_register"     // message type for act register
	MT_PrivateChat = "mt_private_chat" // message type for act send rivate chat message
	MT_GroupChat   = "mt_group_chat"   // message type for act send broadcast chat message
	MT_CreateGroup = "mt_create_group" // message type for act create new group
	MT_JoinGroup   = "mt_join_group"   // message type for act join a group
	MT_LeaveGroup  = "mt_leave_group"  // message type for act leave from group
)

type Message struct {
	MessageType string `json:"message_type"`
	Content     any    `json:"content"`
}

type Sender struct {
	Name     string
	Username string
}

type ClientMessage struct {
	RoomID      uuid.UUID
	RoomName    string
	IsError     bool
	MessageType string
	Target      []string
	Sender      Sender
	Message     []byte
}
