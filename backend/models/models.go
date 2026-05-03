package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleSuperAdmin Role = "super_admin"
)

type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password    string         `json:"-" gorm:"size:255;not null"`
	Email       string         `json:"email" gorm:"size:100"`
	Phone       string         `json:"phone" gorm:"size:20"`
	RealName    string         `json:"real_name" gorm:"size:50"`
	Role        Role           `json:"role" gorm:"size:20;default:user"`
	Status      int            `json:"status" gorm:"default:1"`
	DepartmentID *uint         `json:"department_id"`
	Department  *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type Department struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	ParentID  *uint          `json:"parent_id"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DocumentType struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Code      string         `json:"code" gorm:"size:50;uniqueIndex"`
	ParentID  *uint          `json:"parent_id"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Document struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Title          string         `json:"title" gorm:"size:200;not null"`
	Content        string         `json:"content" gorm:"type:text"`
	DocumentNumber string        `json:"document_number" gorm:"size:100"`
	DocumentTypeID uint           `json:"document_type_id"`
	DocumentType   *DocumentType  `json:"document_type,omitempty" gorm:"foreignKey:DocumentTypeID"`
	DepartmentID   *uint          `json:"department_id"`
	Department     *Department    `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	FileURL        string         `json:"file_url" gorm:"size:500"`
	Author         string         `json:"author" gorm:"size:50"`
	Keyword        string         `json:"keyword" gorm:"size:200"`
	Status         int            `json:"status" gorm:"default:1"`
	ViewCount      int            `json:"view_count" gorm:"default:0"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type Favorite struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"uniqueIndex:idx_user_document"`
	DocumentID uint           `json:"document_id" gorm:"uniqueIndex:idx_user_document"`
	Document   *Document      `json:"document,omitempty" gorm:"foreignKey:DocumentID"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type News struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text"`
	Author    string         `json:"author" gorm:"size:50"`
	IsTop     int            `json:"is_top" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	ViewCount int            `json:"view_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	ImageURL  string         `json:"image_url" gorm:"size:500;not null"`
	LinkURL   string         `json:"link_url" gorm:"size:500"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"size:200;not null"`
	Content   string         `json:"content" gorm:"type:text"`
	IsTop     int            `json:"is_top" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
