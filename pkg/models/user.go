package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	StatusOnline    UserStatus = "online"
	StatusOffline   UserStatus = "offline"
	StatusAway      UserStatus = "away"
	StatusBusy      UserStatus = "busy"
	StatusInvisible UserStatus = "invisible"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"-" gorm:"not null"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Avatar    string     `json:"avatar"`
	Status    UserStatus `json:"status" gorm:"default:'offline'"`
	LastSeen  *time.Time `json:"last_seen"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	SentMessages     []Message `json:"-" gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `json:"-" gorm:"foreignKey:ReceiverID"`
	Contacts         []Contact `json:"-" gorm:"foreignKey:UserID"`
	ChatMembers      []ChatMember `json:"-" gorm:"foreignKey:UserID"`
	InitiatedCalls   []Call `json:"-" gorm:"foreignKey:CallerID"`
	ReceivedCalls    []Call `json:"-" gorm:"foreignKey:CalleeID"`
}

type Contact struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ContactID uuid.UUID `json:"contact_id" gorm:"type:uuid;not null"`
	Nickname  string    `json:"nickname"`
	IsBlocked bool      `json:"is_blocked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User    User `json:"user" gorm:"foreignKey:UserID"`
	Contact User `json:"contact" gorm:"foreignKey:ContactID"`
}

type UserSession struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Token     string    `json:"token" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}