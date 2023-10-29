package ds

import (
	// "gorm.io/gorm"
	"time"
)

const DRAFT string = "черновик"
const FORMED string = "сформирован"
const COMPELTED string = "завершён"
const REJECTED string = "отклонён"
const DELETED string = "удалён"

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"  json:"-"`
	Login     string `gorm:"size:30;not null"  json:"-"`
	Password  string `gorm:"size:30;not null"  json:"-"`
	Name      string `gorm:"size:50;not null"  json:"name"`
	Moderator bool   `gorm:"not null"  json:"-"`
}

type Recipient struct {
	UUID      string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	FIO       string  `gorm:"size:100;not null" json:"fio" binding:"required"`
	ImageURL  *string `gorm:"size:100;not null" json:"image_url" binding:"required"`
	Email     string  `gorm:"size:75;not null" json:"email" binding:"required"`
	Age       int     `gorm:"not null" json:"age" binding:"required"`
	Adress    string  `gorm:"size:100;not null" json:"adress" binding:"required"`
	IsDeleted bool    `gorm:"not null;default:false" json:"-" binding:"-"`
}

type Notification struct {
	UUID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status           string     `gorm:"size:20;not null"`
	CreationDate     time.Time  `gorm:"not null;type:date"`
	FormationDate    *time.Time `gorm:"type:date"`
	CompletionDate   *time.Time `gorm:"type:date"`
	ModeratorId      *string    `json:"-"`
	CustomerId       string     `gorm:"not null"`
	NotificationType string     `gorm:"size:50;not null"`

	Moderator *User 
	Customer  User 
}

type NotificationContent struct {
	RecipientId    string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"recipient_id"`
	NotificationId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"notification_id"`

	Recipient    *Recipient    `gorm:"foreignKey:RecipientId" json:"recipient"`
	Notification *Notification `gorm:"foreignKey:NotificationId" json:"notification"`
}
