package delivery

import (
	"net/http"
	"strconv"

	"rest-apishka/internal/auth"
	"rest-apishka/internal/model"

	"github.com/gin-gonic/gin"
)

// GetFeedbacks godoc
// @Summary Получение списка доставок
// @Description Возвращает список всех не удаленных доставок
// @Tags Доставка
// @Produce json
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param feedbackStatus query string false "Статус доставки" Format(email)
// @Success 200 {object} model.FeedbackRequest "Список багажей"
// @Failure 500 {object} model.FeedbackRequest "Ошибка сервера"
// @Router /feedback [get]
func (h *Handler) GetFeedbacks(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	feedbackStatus := c.DefaultQuery("feedbackStatus", "")

	var feedbacks []model.FeedbackRequest
	var err error

	if authInstance.Role == "moderator" {
		feedbacks, err = h.UseCase.GetFeedbacksModerator(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	} else {
		feedbacks, err = h.UseCase.GetFeedbacksUser(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// GetFeedbackByID godoc
// @Summary Получение доставки по идентификатору
// @Description Возвращает информацию о доставке по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id} [get]
func (h *Handler) GetFeedbackByID(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
		return
	}

	var feedback model.FeedbackGetResponse

	if authInstance.Role == "moderator" {
		feedback, err = h.UseCase.GetFeedbackByIDModerator(uint(feedbackID), authInstance.UserID)
	} else {
		// Получение доставки для пользователя
		feedback, err = h.UseCase.GetFeedbackByIDUser(uint(feedbackID), authInstance.UserID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// DeleteFeedback godoc
// @Summary Удаление доставки
// @Description Удаляет доставку по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param searchFlightNumber query string false "Номер рейса" Format(email)
// @Param startFormationDate query string false "Начало даты формирования" Format(email)
// @Param endFormationDate query string false "Конец даты формирования" Format(email)
// @Param feedbackStatus query string false "Статус доставки" Format(email)
// @Success 200 {object} model.FeedbackRequest "Список багажей"
// @Failure 400 {object} model.FeedbackRequest "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.FeedbackRequest "Ошибка сервера"
// @Router /feedback/{id}/delete [delete]
func (h *Handler) DeleteFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchFlightNumber := c.DefaultQuery("searchFlightNumber", "")
	startFormationDate := c.DefaultQuery("startFormationDate", "")
	endFormationDate := c.DefaultQuery("endFormationDate", "")
	feedbackStatus := c.DefaultQuery("feedbackStatus", "")
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
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

	feedbacks, err := h.UseCase.GetFeedbacksUser(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"feedbacks": feedbacks})
}

// UpdateFeedbackFlightNumber godoc
// @Summary Обновление номера рейса доставки
// @Description Обновляет номер рейса для доставки по её идентификатору
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param flightNumber body model.FeedbackUpdateFlightNumberRequest true "Новый номер рейса"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор доставки или ошибка чтения JSON объекта"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id}/update [put]
func (h *Handler) UpdateFeedbackFlightNumber(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
		return
	}

	var flightNumber model.FeedbackUpdateFlightNumberRequest
	if err := c.BindJSON(&flightNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ошибка чтения JSON объекта"})
		return
	}

	if authInstance.Role == "moderator" {
		err = h.UseCase.UpdateFlightNumberModerator(uint(feedbackID), authInstance.UserID, flightNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feedback, err := h.UseCase.GetFeedbackByIDModerator(uint(feedbackID), authInstance.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"feedback": feedback})
	} else {
		err = h.UseCase.UpdateFlightNumberUser(uint(feedbackID), authInstance.UserID, flightNumber)
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
	}
}

// UpdateFeedbackStatusUser godoc
// @Summary Обновление статуса доставки для пользователя
// @Description Обновляет статус доставки для пользователя по идентификатору доставки
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор доставки"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id}/user [put]
func (h *Handler) UpdateFeedbackStatusUser(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недоупстимый ИД доставки"})
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
// @Summary Обновление статуса доставки для модератора
// @Description Обновляет статус доставки для модератора по идентификатору доставки
// @Tags Доставка
// @Produce json
// @Param id path int true "Идентификатор доставки"
// @Param feedbackStatus body model.FeedbackUpdateStatusRequest true "Новый статус доставки"
// @Success 200 {object} model.FeedbackGetResponse "Информация о доставке"
// @Failure 400 {object} model.FeedbackGetResponse "Недопустимый идентификатор доставки или ошибка чтения JSON объекта"
// @Failure 500 {object} model.FeedbackGetResponse "Ошибка сервера"
// @Router /feedback/{id}/status [put]
func (h *Handler) UpdateFeedbackStatusModerator(c *gin.Context) {
	// Получение экземпляра singleton для аутентификации
	authInstance := auth.GetAuthInstance()

	feedbackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД доставки"})
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
