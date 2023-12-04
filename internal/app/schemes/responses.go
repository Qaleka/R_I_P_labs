package schemes

import (
	"R_I_P_labs/internal/app/ds"
	"time"
)

type AllRecipientsResponse struct {
	Recipients []ds.Recipient `json:"recipients"`
}

type NotificationShort struct {
	UUID           string `json:"uuid"`
	RecipientCount int    `json:"recipient_count"`
}

type GetAllRecipientsResponse struct {
	DraftNotification *NotificationShort         `json:"draft_notification"`
	Recipients            []ds.Recipient `json:"recipients"`
}

type AllNotificationsResponse struct {
	Notifications []NotificationOutput `json:"notifications"`
}

type NotificationResponse struct {
	Notification NotificationOutput `json:"notification"`
	Recipients    []ds.Recipient  `json:"recipients"`
}

type UpdateNotificationResponse struct {
	Notification NotificationOutput  `json:"notifications"`
}

type NotificationOutput struct {
	UUID           string  `json:"uuid"`
	Status         string  `json:"status"`
	CreationDate   string  `json:"creation_date"`
	FormationDate  *string `json:"formation_date"`
	CompletionDate *string `json:"completion_date"`
	Moderator      *string `json:"moderator"`
	Customer       string  `json:"customer"`
	NotificationType      string  `json:"notification_type"`
}

func ConvertNotification(notification *ds.Notification) NotificationOutput {
	output := NotificationOutput{
		UUID:         notification.UUID,
		Status:       notification.Status,
		CreationDate: notification.CreationDate.Format("2006-01-02 15:04:05"),
		NotificationType:    notification.NotificationType,
		Customer:     notification.Customer.Login,
	}

	if notification.FormationDate != nil {
		formationDate := notification.FormationDate.Format("2006-01-02 15:04:05")
		output.FormationDate = &formationDate
	}

	if notification.CompletionDate != nil {
		completionDate := notification.CompletionDate.Format("2006-01-02 15:04:05")
		output.CompletionDate = &completionDate
	}

	if notification.Moderator != nil {
		output.Moderator = &notification.Moderator.Login
	}

	return output
}

type LoginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
}

type SwaggerLoginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type RegisterResp struct {
	Ok bool `json:"ok"`
}