package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CallType string

const (
	CallTypeVoice CallType = "voice"
	CallTypeVideo CallType = "video"
)

type CallStatus string

const (
	CallStatusInitiated CallStatus = "initiated"
	CallStatusRinging   CallStatus = "ringing"
	CallStatusActive    CallStatus = "active"
	CallStatusEnded     CallStatus = "ended"
	CallStatusMissed    CallStatus = "missed"
	CallStatusDeclined  CallStatus = "declined"
	CallStatusFailed    CallStatus = "failed"
)

type Call struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CallerID     uuid.UUID  `json:"caller_id" gorm:"type:uuid;not null"`
	CalleeID     uuid.UUID  `json:"callee_id" gorm:"type:uuid;not null"`
	ChatID       *uuid.UUID `json:"chat_id" gorm:"type:uuid"`
	Type         CallType   `json:"type" gorm:"default:'voice'"`
	Status       CallStatus `json:"status" gorm:"default:'initiated'"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Duration     int        `json:"duration"` // in seconds
	IsScreenShare bool      `json:"is_screen_share" gorm:"default:false"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Caller      User            `json:"caller" gorm:"foreignKey:CallerID"`
	Callee      User            `json:"callee" gorm:"foreignKey:CalleeID"`
	Chat        *Chat           `json:"chat" gorm:"foreignKey:ChatID"`
	Participants []CallParticipant `json:"participants" gorm:"foreignKey:CallID"`
}

type CallParticipant struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CallID    uuid.UUID  `json:"call_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	JoinedAt  time.Time  `json:"joined_at"`
	LeftAt    *time.Time `json:"left_at"`
	IsMuted   bool       `json:"is_muted" gorm:"default:false"`
	IsVideoOn bool       `json:"is_video_on" gorm:"default:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Relationships
	Call Call `json:"call" gorm:"foreignKey:CallID"`
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type CallSettings struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	AllowVoiceCalls   bool      `json:"allow_voice_calls" gorm:"default:true"`
	AllowVideoCalls   bool      `json:"allow_video_calls" gorm:"default:true"`
	AllowScreenShare  bool      `json:"allow_screen_share" gorm:"default:true"`
	AutoAnswer        bool      `json:"auto_answer" gorm:"default:false"`
	RingtonePath      string    `json:"ringtone_path"`
	CallQuality       string    `json:"call_quality" gorm:"default:'auto'"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}