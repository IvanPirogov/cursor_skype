package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeFile  MessageType = "file"
	MessageTypeImage MessageType = "image"
	MessageTypeVideo MessageType = "video"
	MessageTypeAudio MessageType = "audio"
	MessageTypeSystem MessageType = "system"
)

type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

type Message struct {
	ID         uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SenderID   uuid.UUID     `json:"sender_id" gorm:"type:uuid;not null"`
	ReceiverID *uuid.UUID    `json:"receiver_id" gorm:"type:uuid"`
	ChatID     *uuid.UUID    `json:"chat_id" gorm:"type:uuid"`
	Content    string        `json:"content"`
	Type       MessageType   `json:"type" gorm:"default:'text'"`
	Status     MessageStatus `json:"status" gorm:"default:'sent'"`
	IsEdited   bool          `json:"is_edited" gorm:"default:false"`
	ReplyToID  *uuid.UUID    `json:"reply_to_id" gorm:"type:uuid"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Sender    User       `json:"sender" gorm:"foreignKey:SenderID"`
	Receiver  *User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Chat      *Chat      `json:"chat" gorm:"foreignKey:ChatID"`
	ReplyTo   *Message   `json:"reply_to" gorm:"foreignKey:ReplyToID"`
	Replies   []Message  `json:"replies" gorm:"foreignKey:ReplyToID"`
	Files     []File     `json:"files" gorm:"foreignKey:MessageID"`
	Reactions []Reaction `json:"reactions" gorm:"foreignKey:MessageID"`
}

type File struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	FileName  string    `json:"file_name"`
	FileSize  int64     `json:"file_size"`
	MimeType  string    `json:"mime_type"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Message Message `json:"message" gorm:"foreignKey:MessageID"`
}

type Reaction struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Message Message `json:"message" gorm:"foreignKey:MessageID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}

type MessageRead struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MessageID uuid.UUID `json:"message_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ReadAt    time.Time `json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Message Message `json:"message" gorm:"foreignKey:MessageID"`
	User    User    `json:"user" gorm:"foreignKey:UserID"`
}