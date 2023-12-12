package ds

import (
	"R_I_P_labs/internal/app/role"
	"time"
)

const StatusDraft string = "черновик"
const StatusFormed string = "сформирован"
const StatusCompleted string = "завершён"
const StatusRejected string = "отклонён"
const StatusDeleted string = "удалён"

const SendingCompleted string = "отправлено"
const SendingFailed string = "отправка отменена"
const SendingStarted string = "отправка начата"

type User struct {
	UUID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Role     role.Role
	Login    string `gorm:"size:30;not null" json:"login"`
	Password string `gorm:"size:40;not null" json:"-"`
	// The SHA-1 hash is 20 bytes. When encoded in hexadecimal, each byte is represented by two characters. Therefore, the resulting hash string will be 40 characters long
}

type Recipient struct {
	UUID      string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	FIO       string  `gorm:"size:100;not null" form:"fio" json:"fio" binding:"required"`
	ImageURL  *string `gorm:"size:100" json:"image_url" binding:"-"`
	Email     string  `gorm:"size:75;not null" form:"email" json:"email" binding:"required"`
	Age       int     `gorm:"not null" json:"age" form:"age" binding:"required"`
	Adress    string  `gorm:"size:100;not null" form:"adress" json:"adress" binding:"required"`
	IsDeleted bool    `gorm:"not null;default:false" json:"-" binding:"-"`
}

type Notification struct {
	UUID             string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status           string     `gorm:"size:20;not null"`
	CreationDate     time.Time  `gorm:"not null;type:timestamp"`
	FormationDate    *time.Time `gorm:"type:timestamp"`
	CompletionDate   *time.Time `gorm:"type:timestamp"`
	ModeratorId      *string    `json:"-"`
	CustomerId       string     `gorm:"not null"`
	NotificationType *string     `gorm:"size:50"`
	SendingStatus *string    `gorm:"size:40"`

	Moderator *User
	Customer  User
}

type NotificationContent struct {
	RecipientId    string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"recipient_id"`
	NotificationId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"notification_id"`

	Recipient    *Recipient    `gorm:"foreignKey:RecipientId" json:"recipient"`
	Notification *Notification `gorm:"foreignKey:NotificationId" json:"notification"`
}
