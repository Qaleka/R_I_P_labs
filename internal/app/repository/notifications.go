package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

func (r *Repository) GetAllNotifications(formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Notification, error) {
	var notifications []ds.Notification
	var err error

	if formationDateStart != nil && formationDateEnd != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd).
			Find(&notifications).Error
	} else if formationDateStart != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("formation_date >= ?", *formationDateStart).
			Find(&notifications).Error
	} else if formationDateEnd != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("formation_date <= ?", *formationDateEnd).
			Find(&notifications).Error
	} else {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Find(&notifications).Error
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

func (r *Repository) CreateDraftNotification(customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{CreationDate: time.Now(), CustomerId: customerId, Status: ds.DRAFT}
	err := r.db.Create(notification).Error
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *Repository) GetNotificationById(notificationId, customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{}
	err := r.db.Preload("Moderator").Preload("Customer").
	First(notification, ds.Notification{UUID: notificationId, CustomerId: customerId}).Error
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

func (r *Repository) DeleteFromNotification(notificationId, recipientId string) error {
	err := r.db.Delete(&ds.NotificationContent{NotificationId: notificationId, RecipientId: recipientId}).Error
	if err != nil {
		return err
	}
	return nil
}
