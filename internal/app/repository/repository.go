package repository

import (
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetRecipientByID(id string) (*ds.Recipient, error) { // ?
	recipient := &ds.Recipient{}
	err := r.db.Where("recipient_id = ?", id).First(recipient).Error
	if err != nil {
		return nil, err
	}

	return recipient, nil
}

func (r *Repository) GetAllRecipients() ([]ds.Recipient, error) { //FIO ?
	var recipients []ds.Recipient

	err := r.db.Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	return recipients, nil
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

func (r *Repository) DeleteRecipient(id string) error {
	err := r.db.Exec("UPDATE recipients SET is_deleted = ? WHERE recipient_id = ?", true, id).Error
	if err != nil {
		return err
	}

	return nil
}
