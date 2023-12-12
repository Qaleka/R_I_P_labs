package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

func (r *Repository) GetAllNotifications(customerId *string, formationDateStart, formationDateEnd *time.Time, status string) ([]ds.Notification, error) {
	var notifications []ds.Notification
	query := r.db.Preload("Customer").Preload("Moderator").
		Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
		Where("status != ? AND status != ?", ds.StatusDeleted, ds.StatusDraft)

	if customerId != nil {
		query = query.Where("customer_id = ?", *customerId)
	}

	if formationDateStart != nil && formationDateEnd != nil {
		query = query.Where("formation_date BETWEEN ? AND ?", *formationDateStart, *formationDateEnd)
	} else if formationDateStart != nil {
		query = query.Where("formation_date >= ?", *formationDateStart)
	} else if formationDateEnd != nil {
		query = query.Where("formation_date <= ?", *formationDateEnd)
	}
	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (r *Repository) GetDraftNotification(customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{}
	err := r.db.First(notification, ds.Notification{Status: ds.StatusDraft, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return notification, nil
}

func (r *Repository) CreateDraftNotification(customerId string) (*ds.Notification, error) {
	notification := &ds.Notification{CreationDate: time.Now(), CustomerId: customerId, Status: ds.StatusDraft}
	err := r.db.Create(notification).Error
	if err != nil {
		return nil, err
	}
	return notification, nil
}

func (r *Repository) GetNotificationById(notificationId string, userId  *string) (*ds.Notification, error) {
	notification := &ds.Notification{}
	query := r.db.Preload("Moderator").Preload("Customer").
		Where("status != ?", ds.StatusDeleted)
	if userId != nil {
		query = query.Where("customer_id = ?", userId)
	}
	err := query.First(notification, ds.Notification{UUID: notificationId}).Error
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

	err := r.db.Table("notification_contents").
		Select("recipients.*").
		Joins("JOIN recipients ON notification_contents.recipient_id = recipients.uuid").
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

func (r *Repository) DeleteFromNotification(notificationId, RecipientId string) error {
	err := r.db.Delete(&ds.NotificationContent{NotificationId: notificationId, RecipientId: RecipientId}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) CountRecipients(notificationId string) (int64, error) {
	var count int64
	err := r.db.Model(&ds.NotificationContent{}).
		Where("notification_id = ?", notificationId).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}