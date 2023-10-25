package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

func (r *Repository) GetAllNotifications(formationDate *time.Time, status string) ([]ds.Notification, error) {
	var notifications []ds.Notification
	var err error

	if formationDate != nil {
		err = r.db.Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("formation_date = ?", formationDate).
			Find(&notifications).Error
	} else if status != "" {
		err = r.db.Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Find(&notifications).Error
	} else {
		err = r.db.Find(&notifications).Error
	}
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *Repository) GetDraftNotification(customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{}
	err := r.db.First(notification, ds.Notification{Status: ds.DRAFT, CustomerId: customerId}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		notification = &ds.Notification{CreationDate: time.Now(), CustomerId: customerId, Status: ds.DRAFT}
		err := r.db.Create(notification).Error
		if err != nil {
			return nil, err
		}
	}
	return notification, nil
}

func (r *Repository) GetNotificationById(notificationId, customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{}
	err := r.db.First(notification, ds.Notification{UUID: notificationId, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return notification, nil
}

func (r *Repository) GetNotificationContent(notificationId string) ([]ds.Recipient, error) {
	var recipients []ds.Recipient

	err := r.db.Table("notification_content").
		Select("recipients.*").
		Joins("JOIN recipients ON notification_content.recipient_id = recipients.uuid").
		Where(ds.NotificationContent{NotificationId: notificationId}).
		Scan(&recipients).Error

	if err != nil {
		return nil, err
	}
	return recipients, nil
}

func (r *Repository) SaveNotification(notification *ds.Notification) error {
	err := r.db.Save(notification).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFromTransportation(notificationId, recipientId string) error {
	err := r.db.Delete(&ds.NotificationContent{NotificationId: notificationId, RecipientId: recipientId}).Error
	if err != nil {
		return err
	}
	return nil
}
