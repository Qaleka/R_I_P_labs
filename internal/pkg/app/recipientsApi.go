package app

import (
	"fmt"
	"net/http"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/schemes"


	"github.com/gin-gonic/gin"
)

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
	recipient.ImageURL = nil
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