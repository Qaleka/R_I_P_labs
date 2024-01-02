package schemes

import (
	"R_I_P_labs/internal/app/ds"
	"fmt"
)

type AllRecipientsResponse struct {
	Recipients []ds.Recipient `json:"recipients"`
}

type GetAllRecipientsResponse struct {
	DraftNotification *string         `json:"draft_notification"`
	Recipients        []ds.Recipient  `json:"recipients"`
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
	NotificationType      *string  `json:"notification_type"`
	SendingStatus *string `json:"sending_status"`
}

func ConvertNotification(notification *ds.Notification) NotificationOutput {
	output := NotificationOutput{
		UUID:         notification.UUID,
		Status:       notification.Status,
		CreationDate: notification.CreationDate.Format("2006-01-02T15:04:05Z07:00"),
		NotificationType:    notification.NotificationType,
		SendingStatus: notification.SendingStatus,
		Customer:     notification.Customer.Login,
	}

	if notification.FormationDate != nil {
		formationDate := notification.FormationDate.Format("2006-01-02T15:04:05Z07:00")
		output.FormationDate = &formationDate
	}

	if notification.CompletionDate != nil {
		completionDate := notification.CompletionDate.Format("2006-01-02T15:04:05Z07:00")
		output.CompletionDate = &completionDate
	}

	if notification.Moderator != nil {
		fmt.Println(notification.Moderator.Login)
		output.Moderator = &notification.Moderator.Login
		fmt.Println(*output.Moderator)
	}

	return output
}

type AddToNotificationResp struct {
	RecipientsCount int64 `json:"recipient_count"`
}

type AuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

