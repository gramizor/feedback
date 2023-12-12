package delivery

import (
	"net/http"
	"strconv"

	"rest-apishka/internal/auth"
	"rest-apishka/internal/model"

	"github.com/gin-gonic/gin"
)

// GetFeedbacks godoc
// @Summary Получение списка опросов
// @Description Возвращает список всех не удаленных опросов
// @Tags Опрос
// @Produce json
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param feedbackStatus query string false "Статус опроса" Format(email)
// @Success 200 {object} model.FeedbackRequest "Список групп"
// @Failure 500 {object} model.FeedbackRequest "Ошибка сервера"
// @Router /feedback [get]

func (h *Handler) GetFeedbacks(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	feedbackStatus := c.DefaultQuery("feedbackStatus", "")

	var feedbacks []model.FeedbackRequest
	var err error

	if authInstance.Role == "moderator" {
		feedbacks, err = h.UseCase.GetFeedbacksModerator(startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	} else {
		feedbacks, err = h.UseCase.GetFeedbacksUser(startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// GetFeedbackByID godoc
// @Summary Получение опроса по идентификатору
// @Description Возвращает информацию об опросе по её идентификатору
// @Tags Опрос
// @Produce json
// @Param id path int true "Идентификатор опроса"
// @Success 200 {object} model.FeedbackGetResponse "Информация об опросе"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор опроса"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id} [get]
func (h *Handler) GetFeedbackByID(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД опроса"})
		return
	}

	var feedback model.FeedbackGetResponse

	if authInstance.Role == "moderator" {
		feedback, err = h.UseCase.GetFeedbackByIDModerator(uint(feedbackID), authInstance.UserID)
	} else {
		// Получение опроса для пользователя
		feedback, err = h.UseCase.GetFeedbackByIDUser(uint(feedbackID), authInstance.UserID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// DeleteFeedback godoc
// @Summary Удаление опроса
// @Description Удаляет опрос по её идентификатору
// @Tags Опрос
// @Produce json
// @Param id path int true "Идентификатор опроса"
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param feedbackStatus query string false "Статус опроса" Format(email)
// @Success 200 {object} model.FeedbackRequest "Список групп"
// @Failure 400 {object} model.FeedbackRequest "Недопустимый идентификатор опроса"
// @Failure 500 {object} model.FeedbackRequest "Ошибка сервера"
// @Router /feedback/{id}/delete [delete]
func (h *Handler) DeleteFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	feedbackStatus := c.DefaultQuery("feedbackStatus", "")
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД опроса"})
		return
	}

	if authInstance.Role == "moderator" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос недоступен для модератора"})
		return
	}

	err = h.UseCase.DeleteFeedbackUser(uint(feedbackID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	feedbacks, err := h.UseCase.GetFeedbacksUser(startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// UpdateFeedbackStatusUser godoc
// @Summary Обновление статуса опроса для пользователя
// @Description Обновляет статус опроса для пользователя по идентификатору опроса
// @Tags Опрос
// @Produce json
// @Param id path int true "Идентификатор опроса"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор опроса"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id}/user [put]
func (h *Handler) UpdateFeedbackStatusUser(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД опроса"})
		return
	}

	if authInstance.Role == "user" {
		err = h.UseCase.UpdateFeedbackStatusUser(uint(feedbackID), authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feedback, err := h.UseCase.GetFeedbackByIDUser(uint(feedbackID), authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"feedback": feedback})
	} else if authInstance.Role == "moderator" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только пользователю"})
		return
	}
}

// UpdateFeedbackStatusModerator godoc
// @Summary Обновление статуса опроса для модератора
// @Description Обновляет статус опроса для модератора по идентификатору опроса
// @Tags Опрос
// @Produce json
// @Param id path int true "Идентификатор опроса"
// @Param feedbackStatus body model.FeedbackUpdateStatusRequest true "Новый статус опроса"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор опроса или ошибка чтения JSON объекта"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id}/status [put]
func (h *Handler) UpdateFeedbackStatusModerator(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД опроса"})
		return
	}

	var feedbackStatus model.FeedbackUpdateStatusRequest
	if err := c.BindJSON(&feedbackStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authInstance.Role == "moderator" {
		err = h.UseCase.UpdateFeedbackStatusModerator(uint(feedbackID), authInstance.UserID, feedbackStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feedback, err := h.UseCase.GetFeedbackByIDUser(uint(feedbackID), authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"feedback": feedback})
	} else if authInstance.Role == "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "данный запрос доступен только модератору"})
		return
	}
}
