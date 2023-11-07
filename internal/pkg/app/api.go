package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (app *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension

	_, err = app.minioClient.PutObject(c, app.config.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", app.config.MinioEndpoint, app.config.BucketName, imageName)
	return &imageURL, nil
}

// вопрос
func (app *Application) getCustomer() string {
	return "2d217868-ab6d-41fe-9b34-7809083a2e8a"
}

func (app *Application) getModerator() *string {
	moderaorId := "87d54d58-1e24-4cca-9c83-bd2523902729"
	return &moderaorId
}

func (app *Application) GetAllRecipients(c *gin.Context) {
	var request schemes.GetAllRecipientsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipients, err := app.repo.GetRecipientByName(request.FIO)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	draftNotification, err := app.repo.GetDraftNotification(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response := schemes.GetAllRecipientsResponse{DraftNotification: nil, Recipients: recipients}
	if draftNotification != nil {
		response.DraftNotification = &schemes.NotificationShort{UUID: draftNotification.UUID}
		containers, err := app.repo.GetNotificationContent(draftNotification.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		response.DraftNotification.RecipientCount = len(containers)
	}
	c.JSON(http.StatusOK, response)
}

func (app *Application) GetRecipient(c *gin.Context) {
	var request schemes.RecipientRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetRecipientByID(request.RecipientId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}
	c.JSON(http.StatusOK, recipient)
}

func (app *Application) DeleteRecipient(c *gin.Context) {
	var request schemes.RecipientRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetRecipientByID(request.RecipientId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}
	recipient.IsDeleted = true
	if err := app.repo.SaveRecipient(recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) AddRecipient(c *gin.Context) {
	var request schemes.AddRecipientRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient := ds.Recipient(request.Recipient)
	if err := app.repo.AddRecipient(&recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, recipient.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		recipient.ImageURL = imageURL
	}
	if err := app.repo.SaveRecipient(&recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (app *Application) ChangeRecipient(c *gin.Context) {
	var request schemes.ChangeRecipientRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	recipient, err := app.repo.GetRecipientByID(request.RecipientId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}

	if request.FIO != nil {
		recipient.FIO = *request.FIO
	}
	if request.Image != nil {
		imageURL, err := app.uploadImage(c, request.Image, recipient.UUID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		recipient.ImageURL = imageURL
	}
	if request.Email != nil {
		recipient.Email = *request.Email
	}
	if request.Age != nil {
		recipient.Age = *request.Age
	}
	if request.Adress != nil {
		recipient.Adress = *request.Adress
	}

	if err := app.repo.SaveRecipient(recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, recipient)
}

func (app *Application) AddToNotification(c *gin.Context) {
	var request schemes.AddToNotificationRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	var err error

	recipient, err := app.repo.GetRecipientByID(request.RecipientId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if recipient == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("получатель не найден"))
		return
	}

	var notification *ds.Notification
	notification, err = app.repo.GetDraftNotification(app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		notification, err = app.repo.CreateDraftNotification(app.getCustomer())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToNotification(notification.UUID, request.RecipientId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var recipients []ds.Recipient
	recipients, err = app.repo.GetNotificationContent(notification.UUID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllRecipientsResponse{Recipients: recipients})
}

func (app *Application) GetAllNotifications(c *gin.Context) {
	var request schemes.GetAllNotificationsRequst
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notifications, err := app.repo.GetAllNotifications(request.FormationDateStart, request.FormationDateEnd, request.Status)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	outputNotifications := make([]schemes.NotificationOutput, len(notifications))
	for i, notification := range notifications {
		outputNotifications[i] = schemes.ConvertNotification(&notification)
	}
	c.JSON(http.StatusOK, schemes.AllNotificationsResponse{Notifications: outputNotifications})
}

func (app *Application) GetNotification(c *gin.Context) {
	var request schemes.NotificationRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notification, err := app.repo.GetNotificationById(request.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	recipients, err := app.repo.GetNotificationContent(request.NotificationId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, schemes.NotificationResponse{Notification: schemes.ConvertNotification(notification), Recipients: recipients})
}

func (app *Application) UpdateNotification(c *gin.Context) {
	var request schemes.UpdateNotificationRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	notification, err := app.repo.GetNotificationById(request.URI.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	notification.NotificationType = request.NotificationType
	if app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.UpdateNotificationResponse{Notification:schemes.ConvertNotification(notification)})
}

func (app *Application) DeleteNotification(c *gin.Context) {
	var request schemes.NotificationRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notification, err := app.repo.GetNotificationById(request.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("увдомление не найдено"))
		return
	}
	notification.Status = ds.DELETED

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) DeleteFromNotification(c *gin.Context) {
	var request schemes.DeleteFromNotificationRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	notification, err := app.repo.GetNotificationById(request.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if notification.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя редактировать уведомление со статусом: %s", notification.Status))
		return
	}

	if err := app.repo.DeleteFromNotification(request.NotificationId, request.RecipientId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	recipients, err := app.repo.GetNotificationContent(request.NotificationId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, schemes.AllRecipientsResponse{Recipients: recipients})
}

func (app *Application) UserConfirm(c *gin.Context) {
	var request schemes.UserConfirmRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	notification, err := app.repo.GetNotificationById(request.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if notification.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать уведомление со статусом %s", notification.Status))
		return
	}
	notification.Status = ds.FORMED
	now := time.Now()
	notification.FormationDate = &now

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) ModeratorConfirm(c *gin.Context) {
	var request schemes.ModeratorConfirmRequest
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if request.Status != ds.COMPELTED && request.Status != ds.REJECTED {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("status %s not allowed", request.Status))
		return
	}

	notification, err := app.repo.GetNotificationById(request.URI.NotificationId, app.getCustomer())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if notification.Status != ds.FORMED {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", notification.Status, request.Status))
		return
	}
	notification.Status = request.Status
	notification.ModeratorId = app.getModerator()
	if request.Status == ds.COMPELTED {
		now := time.Now()
		notification.CompletionDate = &now
	}

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
