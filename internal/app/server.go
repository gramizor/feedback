package app

import (
	"fmt"
	"log"

	"rest-apishka/docs"

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
	// Группа запросов для багажа
	GroupGroup := r.Group("/group")
	{
		GroupGroup.GET("/", app.Handler.GetGroups)
		GroupGroup.GET("/:group_id", app.Handler.GetGroupByID)
		GroupGroup.DELETE("/:group_id/delete", app.Handler.DeleteGroup)
		GroupGroup.POST("/create", app.Handler.CreateGroup)
		GroupGroup.PUT("/:group_id/update", app.Handler.UpdateGroup)
		GroupGroup.POST("/:group_id/feedback", app.Handler.AddGroupToFeedback)
		GroupGroup.DELETE("/:group_id/feedback/delete", app.Handler.RemoveGroupFromFeedback)
		GroupGroup.POST("/:group_id/image", app.Handler.AddGroupImage)
	}

	// Группа запросов для доставки
	FeedbackGroup := r.Group("/feedback")
	{
		FeedbackGroup.GET("/", app.Handler.GetFeedbacks)
		FeedbackGroup.GET("/:id", app.Handler.GetFeedbackByID)
		FeedbackGroup.DELETE("/:id/delete", app.Handler.DeleteFeedback)
		FeedbackGroup.PUT("/:id/update", app.Handler.UpdateFeedbackFlightNumber)
		FeedbackGroup.PUT("/:id/status/user", app.Handler.UpdateFeedbackStatusUser)           // Новый маршрут для обновления статуса доставки пользователем
		FeedbackGroup.PUT("/:id/status/moderator", app.Handler.UpdateFeedbackStatusModerator) // Новый маршрут для обновления статуса доставки модератором
	}

	addr := fmt.Sprintf("%s:%d", app.Config.ServiceHost, app.Config.ServicePort)
	r.Run(addr)
	log.Println("Server down")
}
