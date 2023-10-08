package repository

//вопрос владу
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

	// err := r.db.First(container, "container_id = ?", id).Error
	err := r.db.Where("recipient_id = ?", id).First(recipient).Error
	if err != nil {
		return nil, err
	}

	return recipient, nil

	// nContent := &ds.NotificationContent{}

	// err = r.db.Where("notification_id = ?", recipient.RecipientId).First(nContent).Error
	// if err != nil {
	// 	nContent.Cargo = "Отсутствует"
	// 	nContent.Weight = 0
	// }

	// return &ds.RecipientInfo{
	// 	RecipientId: recipient.RecipientId,
	// 	ImageURL:    recipient.ImageURL,
	// 	FIO:         recipient.FIO,
	// 	Email:       recipient.Email,
	// 	Cargo:       tComposition.Cargo,
	// 	Weight:      tComposition.Weight,
	// }, nil
}

// type RecipientInfo struct {
// 	RecipientId string
// 	ImageURL    string
// 	FIO         string
// 	Email       string
// 	Age         int
// 	Adress      string
// }

func (r *Repository) GetAllRecipients() ([]ds.Recipient, error) { //FIO ?
	var recipients []ds.Recipient

	err := r.db.Preload("FIO").Find(&recipients).Error
	if err != nil {
		return nil, err
	}

	return recipients, nil
}

func (r *Repository) GetRecipientByName(FIO string) ([]ds.Recipient, error) {
	var recipients []ds.Recipient

	err := r.db.
		Where("LOWER(recipients.FIO) LIKE ?", "%"+strings.ToLower(FIO)+"%").
		Find(&recipients).Error

	if err != nil {
		return nil, err
	}

	return recipients, nil
}

func (r *Repository) DeliviredNotification(id string) error {
	err := r.db.Exec("UPDATE recipients SET delivired = ? WHERE recipient_id = ?", true, id).Error
	if err != nil {
		return err
	}

	return nil
}
