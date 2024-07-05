package payload

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// tanpa gorm.Model
// =======================================================================
// User model
type User struct {
	UuId          uuid.UUID      `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid" form:"uuid"`
	UuIdAddr      uuid.UUID      `gorm:"column:uuidaddr" json:"uuidaddr" form:"uuidaddr"`
	UuIdAdmission uuid.UUID      `gorm:"column:uuidAdmission" json:"uuidAdmission" form:"uuidAdmission"`
	Gender        string         `gorm:"column:gender" json:"gender" form:"gender"`
	Name          string         `gorm:"column:name" json:"name" form:"name"`
	Noid          string         `gorm:"column:noid" json:"noid" form:"noid"`
	Religion      string         `gorm:"column:religion" json:"religion" form:"religion"`
	Roles         []Role         `gorm:"many2many:user_roles;" json:"roles" form:"roles"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// =======================================================================

// karena []string tidak bisa maka []Role
// Role model
type Role struct {
	UuId uuid.UUID `gorm:"column:uuid;type:uuid;primaryKey;default:gen_random_uuid()" json:"uuid" form:"uuid"`
	Name string    `gorm:"column:name" json:"name" form:"name"`
}

// =======================================================================
