package app

import (
	"fmt"
	"log"

	"rest-apishka/docs"
	"rest-apishka/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// Run запускает приложение.
func (app *Application) Run() {
	r := gin.Default()
	// Это нужно для автоматического создания папки "docs" в вашем проекте
	docs.SwaggerInfo.Title = "Feedback RestAPI"
	docs.SwaggerInfo.Description = "API server for Feedback application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Группа запросов для группы
	GroupGroup := r.Group("/group")
	{
		GroupGroup.GET("/", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetGroups)
		GroupGroup.GET("/paginate", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetGroupsPaged)
		GroupGroup.GET("/:group", middleware.Guest(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.GetGroupByID)
		GroupGroup.DELETE("/:group/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.DeleteGroup)
		GroupGroup.POST("/create", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.CreateGroup)
		GroupGroup.PUT("/:group/update", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.UpdateGroup)
		GroupGroup.POST("/:group/feedback", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddGroupToFeedback)
		GroupGroup.DELETE("/:group/feedback/delete", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.RemoveGroupFromFeedback)
		GroupGroup.POST("/:group/image", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.AddGroupImage)
	}

	// Группа запросов для опроса
	FeedbackGroup := r.Group("/feedback").Use(middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository))
	{
		FeedbackGroup.GET("/", app.Handler.GetFeedbacks)
		FeedbackGroup.GET("/:id", app.Handler.GetFeedbackByID)
		FeedbackGroup.DELETE("/:id/delete", app.Handler.DeleteFeedback)
		FeedbackGroup.PUT("/:id/status/user", app.Handler.UpdateFeedbackStatusUser)           // Новый маршрут для обновления статуса опроса пользователем
		FeedbackGroup.PUT("/:id/status/moderator", app.Handler.UpdateFeedbackStatusModerator) // Новый маршрут для обновления статуса опроса модератором
	}

	UserGroup := r.Group("/user")
	{
		UserGroup.GET("/", app.Handler.GetUserByID)
		UserGroup.POST("/registration", app.Handler.Register)
		UserGroup.POST("/login", app.Handler.Login)
		UserGroup.POST("/logout", middleware.Authenticate(app.Repository.GetRedisClient(), []byte("AccessSecretKey"), app.Repository), app.Handler.Logout)
	}
	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
