package repository

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

func (r *Repository) GetContainerByID(id string) (*ds.Recipient, error) {
	recipient := &ds.Recipient{UUID: id}
	err := r.db.First(recipient, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return recipient, nil
}

func (r *Repository) AddRecipient(recipient *ds.Recipient) error {
	err := r.db.Create(&recipient).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetRecipientByName(FIO string) ([]ds.Recipient, error) {
	var recipients []ds.Recipient

	err := r.db.
		Where("LOWER(recipients.fio) LIKE ?", "%"+strings.ToLower(FIO)+"%").Where("is_deleted = ?", false).
		Find(&recipients).Error

	if err != nil {
		return nil, err
	}

	return recipients, nil
}

// исправить
func (r *Repository) SaveRecipient(recipient *ds.Recipient) error {
	err := r.db.Save(recipient).Error
	if err != nil {
		return err
	}
	return nil
}
