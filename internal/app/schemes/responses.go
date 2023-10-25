package schemes

import (
	"R_I_P_labs/internal/app/ds"
)

type AllRecipientsResponse struct {
	Recipients []ds.Recipient `json:"recipients"`
}

type AllNotificationsResponse struct {
	Notifications []ds.Notification `json:"notifications"`
}

type NotificationResponse struct {
	Notification ds.Notification `json:"notifications"`
	Recipient    []ds.Recipient  `json:"recipients"`
}

type UpdateNotificationResponse struct {
	Notification ds.Notification `json:"notifications"`
}
