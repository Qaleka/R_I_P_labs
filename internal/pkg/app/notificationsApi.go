package app

import (
	"fmt"
	"net/http"
	"time"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

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
	if notification.Status != ds.DRAFT {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя сформировать уведомление со статусом %s", notification.Status))
		return
	}
	if request.Confirm {
		notification.Status = ds.FORMED
		now := time.Now()
		notification.FormationDate = &now
	} else {
		notification.Status = ds.DELETED
	}

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
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", notification.Status,  ds.FORMED))
		return
	}
	if request.Confirm {
		notification.Status = ds.COMPELTED
		now := time.Now()
		notification.CompletionDate = &now
	
	} else {
		notification.Status = ds.REJECTED
	}
	notification.ModeratorId = app.getModerator()
	
	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}
