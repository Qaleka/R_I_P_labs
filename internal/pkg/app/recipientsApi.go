package app

import (
	"fmt"
	"net/http"

	_ "R_I_P_labs/docs"
	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

// @Summary		Получить всех получателей
// @Tags		Получатели
// @Description	Возвращает всех доуступных получателей с опциональной фильтрацией по ФИО
// @Produce		json
// @Param		fio query string false "ФИО для фильтрации"
// @Success		200 {object} schemes.GetAllRecipientsResponse
// @Router		/api/recipients [get]
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
	response := schemes.GetAllRecipientsResponse{DraftNotification: nil, Recipients: recipients}
	if userId, exists := c.Get("userId"); exists {
		draftNotification, err := app.repo.GetDraftNotification(userId.(string))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if draftNotification != nil {
			response.DraftNotification = &draftNotification.UUID
		}
	}
	c.JSON(http.StatusOK, response)
}

// @Summary		Получить одного получателя
// @Tags		Получатели
// @Description	Возвращает более подробную информацию об одном получателе
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200 {object} ds.Recipient
// @Router		/api/recipients/{id} [get]
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

// @Summary		Удалить получателя
// @Tags		Получатели
// @Description	Удаляет получателя по id
// @Param		id path string true "id получателя"
// @Success		200
// @Router		/api/recipients/{id} [delete]
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
	if recipient.ImageURL != nil {
		if err := app.deleteImage(c, recipient.UUID); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	recipient.ImageURL = nil
	recipient.IsDeleted = true
	if err := app.repo.SaveRecipient(recipient); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Добавить получателя
// @Tags		Получатели
// @Description	Добавить нового получателя
// @Accept		mpfd
// @Param     	image formData file false "Изображение получателя"
// @Param     	fio formData string true "ФИО" format:"string" maxLength:100
// @Param     	email formData string true "Почта" format:"string" maxLength:100
// @Param     	age formData int true "Возраст" format:"int"
// @Param     	adress formData string true "Адрес" format:"string" maxLength:100
// @Success		200
// @Router		/api/recipients [post]
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

	c.Status(http.StatusCreated)
}

// @Summary		Изменить получателя
// @Tags		Получатели
// @Description	Изменить данные полей о получателе
// @Accept		mpfd
// @Param		id path string true "Идентификатор получателя" format:"uuid"
// @Param		fio formData string false "ФИО" format:"string" maxLength:100
// @Param		email formData string false "Почта" format:"string" maxLength:100
// @Param		age formData int false "Возраст" format:"int"
// @Param		image formData file false "Изображение получателя"
// @Param		adress formData string false "Адрес" format:"string" maxLength:100
// @Success		200
// @Router		/api/recipients/{id} [put]
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
		if recipient.ImageURL != nil {
			if err := app.deleteImage(c, recipient.UUID); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
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

	c.Status(http.StatusOK)
}

// @Summary		Добавить в уведомление
// @Tags		Получатели
// @Description	Добавить выбранного получателя в черновик уведомления
// @Produce		json
// @Param		id path string true "id получателя"
// @Success		200
// @Router		/api/recipients/{id}/add_to_notification [post]
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
	userId := getUserId(c)
	notification, err = app.repo.GetDraftNotification(userId)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if notification == nil {
		notification, err = app.repo.CreateDraftNotification(userId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err = app.repo.AddToNotification(notification.UUID, request.RecipientId); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}