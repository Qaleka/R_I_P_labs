package schemes

import (
	"R_I_P_labs/internal/app/ds"

	"mime/multipart"
	"time"
)

type RecipientRequest struct {
	RecipientId string `uri:"id" binding:"required,uuid"`
}

type GetAllRecipientsRequest struct {
	FIO string `form:"fio"`
}

type AddRecipientRequest struct {
	ds.Recipient
	Image *multipart.FileHeader `form:"image" json:"image"`
}

type ChangeRecipientRequest struct {
	RecipientId string                `uri:"id" binding:"required,uuid"`
	FIO         *string               `form:"fio" json:"fio" binding:"omitempty,max=100"`
	Email       *string               `form:"email" json:"email" binding:"omitempty,max=75"`
	Age         *int                  `form:"age" json:"age"`
	Image       *multipart.FileHeader `form:"image" json:"image"`
	Adress      *string               `form:"adress" json:"adress" binding:"omitempty,max=100"`
}

type AddToNotificationRequest struct {
	RecipientId string `uri:"id" binding:"required,uuid"`
}

type GetAllNotificationsRequst struct {
	FormationDateStart *time.Time `form:"formation_date_start" json:"formation_date_start" time_format:"2006-01-02 15:04"`
	FormationDateEnd   *time.Time `form:"formation_date_end" json:"formation_date_end" time_format:"2006-01-02 15:04"`
	Status             string     `form:"status"`
}

type NotificationRequest struct {
	NotificationId string `uri:"id" binding:"required,uuid"`
}

type UpdateNotificationRequest struct {
	NotificationType string `form:"notification_type" json:"notification_type" binding:"required,max=50"`
}

type DeleteFromNotificationRequest struct {
	RecipientId    string `uri:"id" binding:"required,uuid"`
}

type ModeratorConfirmRequest struct {
	URI struct {
		NotificationId string `uri:"id" binding:"required,uuid"`
	}
	Confirm *bool `form:"confirm" binding:"required"`
}

type LoginReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type RegisterReq struct {
	Login    string `form:"login" binding:"required,max=30"`
	Password string `form:"password" binding:"required,max=30"`
}

type SendingReq struct {
	URI struct {
		NotificationId string `uri:"id" binding:"required,uuid"`
	}
	SendingStatus *bool `json:"sending_status" form:"sending_status" binding:"required"`
	Token string `json:"token" form:"token" binding:"required"`
}