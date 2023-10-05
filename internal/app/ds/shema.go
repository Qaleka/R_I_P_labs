package ds

import (
	// "gorm.io/gorm"
	"time"
)

type Status struct {
	StatusId uint   `gorm:"primaryKey"`
	Name     string `gorm:"size:50;not null"`
}

type User struct {
	UserId    uint   `gorm:"primaryKey"`
	Login     string `gorm:"size:30;not null"`
	Password  string `gorm:"size:30;not null"`
	Name      string `gorm:"size:50;not null"`
	Moderator bool   `gorm:"not null"`
}

type Recipient struct {
	RecipientId uint   `gorm:"primaryKey;not null;autoIncrement:false"`
	FIO         string `gorm:"size:100;not null"`
	ImageURL    string `gorm:"size:100;not null"`
	Email       string `gorm:"size:50;not null"`
	Age         uint   `gorm:"not null"`
	Adress      string `gorm:"size:100;not null"`
}

type Notification struct {
	NotificationId   uint       `gorm:"primaryKey"`
	StatusId         uint       `gorm:"not null"`
	CreationDate     time.Time  `gorm:"not null;type:date"`
	FormationDate    *time.Time `gorm:"type:date"`
	CompletionDate   *time.Time `gorm:"type:date"`
	ModeratorId      uint       `gorm:"not null"`
	CustomerId       uint       `gorm:"not null"`
	NotificationType string     `gorm:"size:50;not null"`

	Status    Status
	Moderator User `gorm:"foreignKey:ModeratorId"`
	Customer  User `gorm:"foreignKey:CustomerId"`
}

type NotificationContent struct {
	RecipientId    uint   `gorm:"primaryKey;not null;autoIncrement:false"`
	NotificationId uint   `gorm:"primaryKey;not null;autoIncrement:false"`
	MessageContent string `gorm:"size:100;not null"`

	Recipient    Recipient    `gorm:"foreignKey:RecipientId"`
	Notification Notification `gorm:"foreignKey:NotificationId"`
}