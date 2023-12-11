package app

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"R_I_P_labs/internal/app/config"
	"R_I_P_labs/internal/app/dsn"
	"R_I_P_labs/internal/app/redis"
	"R_I_P_labs/internal/app/repository"
	"R_I_P_labs/internal/app/role"


	_ "R_I_P_labs/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "R_I_P_labs/docs"
)

type Application struct {
	repo        *repository.Repository
	minioClient *minio.Client
	config      *config.Config
	redisClient *redis.Client
}

func (app *Application) Run() {
	log.Println("Server start up")

	r := gin.Default()

	r.Use(ErrorHandler())

	// Услуги - получатели
	api := r.Group("/api")
	{
		recipients := api.Group("/recipients")
		{
			recipients.GET("/", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetAllRecipients)                     // Список с поиском
			recipients.GET("/:recipient_id", app.WithAuthCheck(role.NotAuthorized, role.Customer, role.Moderator), app.GetRecipient)            // Одна услуга
			recipients.DELETE("/:recipient_id", app.WithAuthCheck(role.Moderator), app.DeleteRecipient)                         				// Удаление
			recipients.PUT("/:recipient_id", app.WithAuthCheck(role.Moderator), app.ChangeRecipient)                            				// Изменение
			recipients.POST("/", app.WithAuthCheck(role.Moderator), app.AddRecipient)                                           				// Добавление
			recipients.POST("/:recipient_id/add_to_notification", app.WithAuthCheck(role.Customer,role.Moderator), app.AddToNotification) 		// Добавление в заявку
		}

		// Заявки - уведомления
		notifications := api.Group("/notifications")
		{
			notifications.GET("/", app.WithAuthCheck(role.Customer, role.Moderator), app.GetAllNotifications)                                         				  // Список (отфильтровать по дате формирования и статусу)
			notifications.GET("/:notification_id",app.WithAuthCheck(role.Customer, role.Moderator),  app.GetNotification)                             				  // Одна заявка
			notifications.PUT("/:notification_id/update", app.WithAuthCheck(role.Customer, role.Moderator), app.UpdateNotification)                                	  // Изменение (добавление транспорта)
			notifications.DELETE("/:notification_id", app.WithAuthCheck(role.Customer,role.Moderator), app.DeleteNotification)                                      				  // Удаление
			notifications.DELETE("/:notification_id/delete_recipient/:recipient_id", app.WithAuthCheck(role.Customer, role.Moderator), app.DeleteFromNotification) 	  // Изменеие (удаление услуг)
			notifications.PUT("/user_confirm", app.WithAuthCheck(role.Customer, role.Moderator), app.UserConfirm)                                    				  // Сформировать создателем
			notifications.PUT("/:notification_id/moderator_confirm", app.WithAuthCheck(role.Moderator), app.ModeratorConfirm)                         				  // Завершить отклонить модератором
		}

		// Пользователи (авторизация)
		user := api.Group("/user")
		{
			user.POST("/sign_up", app.Register)
			user.POST("/login", app.Login)
			user.POST("/logout", app.Logout)
		}

		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		r.Run(fmt.Sprintf("%s:%d", app.config.ServiceHost, app.config.ServicePort))

		log.Println("Server down")
	}
}


func New() (*Application, error) {
	var err error
	loc, _ := time.LoadLocation("UTC")
	time.Local = loc
	app := Application{}
	app.config, err = config.NewConfig()
	if err != nil {
		return nil, err
	}

	app.repo, err = repository.New(dsn.FromEnv())
	if err != nil {
		return nil, err
	}

	app.minioClient, err = minio.New(app.config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", "", ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	app.redisClient, err = redis.New(app.config.Redis)
	if err != nil {
		return nil, err
	}

	return &app, nil
}