package app

import (
	"fmt"
	"net/http"
	"time"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/role"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

// @Summary		Получить все уведомления
// @Tags		Уведомления
// @Description	Возвращает все уведомления с фильтрацией по статусу и дате формирования
// @Produce		json
// @Param		status query string false "статус уведомления"
// @Param		formation_date_start query string false "начальная дата формирования"
// @Param		formation_date_end query string false "конечная дата формирвания"
// @Success		200 {object} schemes.AllNotificationsResponse
// @Router		/api/notifications [get]
func (app *Application) GetAllNotifications(c *gin.Context) {
	var request schemes.GetAllNotificationsRequst
	var err error
	if err = c.ShouldBindQuery(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var notifications []ds.Notification
	if userRole == role.Customer {
		notifications, err = app.repo.GetAllNotifications(&userId, request.FormationDateStart, request.FormationDateEnd, request.Status)
	} else {
		notifications, err = app.repo.GetAllNotifications(nil, request.FormationDateStart, request.FormationDateEnd, request.Status)
	}
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

// @Summary		Получить одно уведомление
// @Tags		Уведомления
// @Description	Возвращает подробную информацию об уведомлении и его типе
// @Produce		json
// @Param		id path string true "id уведомления"
// @Success		200 {object} schemes.NotificationResponse
// @Router		/api/notifications/{id} [get]
func (app *Application) GetNotification(c *gin.Context) {
	var request schemes.NotificationRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userId := getUserId(c)
	userRole := getUserRole(c)
	var notification *ds.Notification
	if userRole == role.Moderator {
		notification, err = app.repo.GetNotificationById(request.NotificationId, nil)
	} else {
		notification, err = app.repo.GetNotificationById(request.NotificationId, &userId)
	}
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

type SwaggerUpdateNotificationRequest struct {
	NotificationType string `json:"notification_type"`
}

// @Summary		Указать тип уведомления
// @Tags		Уведомления
// @Description	Позволяет изменить тип чернового уведомления и возвращает обновлённые данные
// @Access		json
// @Produce		json
// @Param		notification_type body SwaggerUpdateNotificationRequest true "Тип уведомления"
// @Success		200
// @Router		/api/notifications [put]
func (app *Application) UpdateNotification(c *gin.Context) {
	var request schemes.UpdateNotificationRequest
	var err error
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Получить черновую заявку
	var notification *ds.Notification

	userId := getUserId(c)
	notification, err = app.repo.GetDraftNotification(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	// Добавить тип
	notification.NotificationType = &request.NotificationType
	if app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Удалить черновое уведомление
// @Tags		Уведомления
// @Description	Удаляет черновое уведомление
// @Success		200
// @Router		/api/notifications [delete]
func (app *Application) DeleteNotification(c *gin.Context) {
	var err error
	// Получить черновую заявку
	var notification *ds.Notification
	userId := getUserId(c)
	notification, err = app.repo.GetDraftNotification(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("увдомление не найдено"))
		return
	}

	notification.Status = ds.StatusDeleted

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Удалить получателя из черновово уведомления
// @Tags		Уведомления
// @Description	Удалить получателя из черновово уведомления
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200
// @Router		/api/notifications/delete_recipient/{id} [delete]
func (app *Application) DeleteFromNotification(c *gin.Context) {
	var request schemes.DeleteFromNotificationRequest
	var err error
	if err := c.ShouldBindUri(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// Получить черновую заявку
	var notification *ds.Notification
	userId := getUserId(c)
	notification, err = app.repo.GetDraftNotification(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}

	if err := app.repo.DeleteFromNotification(notification.UUID, request.RecipientId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Сформировать уведомление
// @Tags		Уведомления
// @Description	Сформировать уведомление пользователем
// @Success		200
// @Router		/api/notifications/user_confirm [put]
func (app *Application) UserConfirm(c *gin.Context) {
	userId := getUserId(c)
	notification, err := app.repo.GetDraftNotification(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	
	if err := sendingRequest(notification.UUID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf(`sending service is unavailable: {%s}`, err))
		return
	}

	sendingStatus := ds.SendingStarted
	notification.SendingStatus = &sendingStatus
	notification.Status = ds.StatusFormed
	now := time.Now()
	notification.FormationDate = &now

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	
	c.Status(http.StatusOK)
}

// @Summary		Подтвердить уведомление
// @Tags		Уведомления
// @Description	Подтвердить или отменить уведомление модератором
// @Param		id path string true "id уведомления"
// @Param		confirm body boolean true "подтвердить"
// @Success		200
// @Router		/api/notifications/{id}/moderator_confirm [put]
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

	userId := getUserId(c)
	notification, err := app.repo.GetNotificationById(request.URI.NotificationId,nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	if notification.Status != ds.StatusFormed  {
		c.AbortWithError(http.StatusMethodNotAllowed, fmt.Errorf("нельзя изменить статус с \"%s\" на \"%s\"", notification.Status,  ds.StatusFormed ))
		return
	}


	if *request.Confirm {
		notification.Status = ds.StatusCompleted
	} else {
		notification.Status = ds.StatusRejected
	}
	now := time.Now()
	notification.CompletionDate = &now
	moderator, err := app.repo.GetUserById(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	notification.Moderator = moderator
	
	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

func (app *Application) Sending(c *gin.Context) {
	var request schemes.SendingReq
	if err := c.ShouldBindUri(&request.URI); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	

	if request.Token != app.config.Token {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	notification, err := app.repo.GetNotificationById(request.URI.NotificationId, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		c.AbortWithError(http.StatusNotFound, fmt.Errorf("уведомление не найдено"))
		return
	}
	// if notification.Status != ds.StatusFormed || *notification.SendingStatus != ds.SendingStarted {
	// 	c.AbortWithStatus(http.StatusMethodNotAllowed)
	// 	return
	// }

	var sendingStatus string
	if *request.SendingStatus {
		sendingStatus = ds.SendingCompleted
	} else {
		sendingStatus = ds.SendingFailed
	}
	notification.SendingStatus = &sendingStatus

	if err := app.repo.SaveNotification(notification); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}