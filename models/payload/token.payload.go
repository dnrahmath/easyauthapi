package payload

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// tanpa gorm.Model
// ===================================================================
// Token model
type Token struct {
	UuId        uuid.UUID      `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid" form:"uuid"`
	UuIdUser    uuid.UUID      `gorm:"column:uuiduser" json:"uuiduser" form:"uuiduser"`
	Token       string         `gorm:"column:token" json:"token" form:"token"`
	Type        string         `gorm:"column:type" json:"type" form:"type"`
	ExpiresAt   time.Time      `gorm:"column:expires_at" json:"expires_at" form:"expires_at"`
	Blacklisted bool           `gorm:"column:blacklisted" json:"blacklisted" form:"blacklisted"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

//===================================================================

type PayloadToken struct {
	Sub  string `json:"sub"`
	Exp  int64  `json:"exp"`
	Iat  int64  `json:"iat"`
	User *User  `json:"user"`
	Type string `json:"type"`
}

//==================================

type DecodeToken struct {
	Header    []byte `json:"header"`
	Payload   []byte `json:"payload"`
	Signature string `json:"signature"`
}

//===================================================================
