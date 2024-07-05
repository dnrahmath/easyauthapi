package payload

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// tanpa gorm.Model
// ===================================================================
// Admission model
type Admission struct {
	UuId                uuid.UUID      `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid" form:"uuid"`
	UuIdUser            uuid.UUID      `gorm:"column:uuiduser" json:"uuiduser" form:"uuiduser"`
	Username            string         `gorm:"column:username" json:"username" form:"username"`
	Password            string         `gorm:"column:password" json:"password" form:"password"`
	Email               string         `gorm:"column:email" json:"email" form:"email"`
	EmailVerified       bool           `gorm:"column:emailverified" json:"emailverified" form:"emailverified"`
	PhoneNumber         string         `gorm:"column:phonenumber" json:"phonenumber" form:"phonenumber"`
	PhoneNumberVerified bool           `gorm:"column:phonenumberverified" json:"phonenumberverified" form:"phonenumberverified"`
	CreatedAt           time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index"`
}

//===================================================================

// type untuk validasi sendcode
type TypeVal struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

func (t *TypeVal) DefaultValue() *TypeVal {
	return &TypeVal{
		Username:    "username",
		Email:       "email",
		PhoneNumber: "phonenumber",
	}
}

//===================================================================
