package ds

import (
	// "gorm.io/gorm"
	"time"
)

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid"`
	Login     string `gorm:"size:30;not null" json:"login"`
	Password  string `gorm:"size:30;not null" json:"password"`
	Name      string `gorm:"size:50;not null"  json:"name"`
	Moderator bool   `gorm:"not null" json:"moderator"`
}

type Recipient struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	FIO       string `gorm:"size:100;not null" json:"fio" binding:"required"`
	ImageURL  string `gorm:"size:100;not null" json:"image_url" binding:"required"`
	Email     string `gorm:"size:75;not null" json:"email" binding:"required"`
	Age       int    `gorm:"not null" json:"age" binding:"required"`
	Adress    string `gorm:"size:100;not null" json:"adress" binding:"required"`
	IsDeleted bool   `gorm:"not null;default:false" json:"is_deleted" binding:"-"`
}

type Notification struct {
	UUID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid"`
	Status           string     `gorm:"size:20;not null" json:"status"`
	CreationDate     time.Time  `gorm:"not null;type:date" json:"creation_date"`
	FormationDate    *time.Time `gorm:"not null;type:date" json:"formation_date"`
	CompletionDate   *time.Time `gorm:"not null;type:date" json:"completion_date"`
	ModeratorId      *string    `json:"moderator_id"`
	CustomerId       string     `gorm:"not null" json:"customer_id"`
	NotificationType string     `gorm:"size:50;not null" json:"notification_type"`

	Moderator User `gorm:"foreignKey:ModeratorId" json:"-"`
	Customer  User `gorm:"foreignKey:CustomerId" json:"-"`
}

type NotificationContent struct {
	RecipientId    uint `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"recipient_id"`
	NotificationId uint `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"notification_id"`

	Recipient    *Recipient    `gorm:"foreignKey:RecipientId" json:"recipient"`
	Notification *Notification `gorm:"foreignKey:NotificationId" json:"notification"`
}
