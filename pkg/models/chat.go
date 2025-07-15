package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type Chat struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        ChatType  `json:"type" gorm:"default:'private'"`
	Avatar      string    `json:"avatar"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Creator  User         `json:"creator" gorm:"foreignKey:CreatedBy"`
	Members  []ChatMember `json:"members" gorm:"foreignKey:ChatID"`
	Messages []Message    `json:"messages" gorm:"foreignKey:ChatID"`
}

type ChatMemberRole string

const (
	ChatMemberRoleAdmin     ChatMemberRole = "admin"
	ChatMemberRoleModerator ChatMemberRole = "moderator"
	ChatMemberRoleMember    ChatMemberRole = "member"
)

type ChatMember struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID    uuid.UUID      `json:"chat_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Role      ChatMemberRole `json:"role" gorm:"default:'member'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	JoinedAt  time.Time      `json:"joined_at"`
	LeftAt    *time.Time     `json:"left_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`

	// Relationships
	Chat Chat `json:"chat" gorm:"foreignKey:ChatID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type ChatSettings struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID            uuid.UUID `json:"chat_id" gorm:"type:uuid;not null"`
	AllowInvites      bool      `json:"allow_invites" gorm:"default:true"`
	AllowMembersAdd   bool      `json:"allow_members_add" gorm:"default:true"`
	AllowFileSharing  bool      `json:"allow_file_sharing" gorm:"default:true"`
	AllowVoiceCalls   bool      `json:"allow_voice_calls" gorm:"default:true"`
	AllowVideoCalls   bool      `json:"allow_video_calls" gorm:"default:true"`
	MessageRetention  int       `json:"message_retention" gorm:"default:0"` // days, 0 = forever
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	Chat Chat `json:"chat" gorm:"foreignKey:ChatID"`
}

type ChatInvite struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ChatID    uuid.UUID `json:"chat_id" gorm:"type:uuid;not null"`
	InvitedBy uuid.UUID `json:"invited_by" gorm:"type:uuid;not null"`
	InvitedUser uuid.UUID `json:"invited_user" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt *time.Time `json:"expires_at"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Chat        Chat `json:"chat" gorm:"foreignKey:ChatID"`
	InvitedByUser User `json:"invited_by_user" gorm:"foreignKey:InvitedBy"`
	InvitedUserRecord User `json:"invited_user_record" gorm:"foreignKey:InvitedUser"`
}