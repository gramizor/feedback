package delivery

import (
	"io"
	"net/http"
	"strconv"

	"rest-apishka/internal/auth"
	"rest-apishka/internal/model"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка багажа
// @Description Возращает список всех активных багажей
// @Tags Багаж
// @Produce json
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список багажей"
// @Failure 500 {object} model.GroupsGetResponse "Ошибка сервера"
// @Router /group [get]
func (h *Handler) GetGroups(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchCode := c.DefaultQuery("searchCode", "")

	groups, err := h.UseCase.GetGroups(searchCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Получение багажа по ID
// @Description Возвращает информацию о багаже по его ID
// @Tags Багаж
// @Produce json
// @Param group_id path int true "ID багажа"
// @Success 200 {object} model.Group "Информация о багаже"
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id} [get]
func (h *Handler) GetGroupByID(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

	group, err := h.UseCase.GetGroupByID(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

// @Summary Создание нового багажа
// @Description Создает новый багаж с предоставленными данными
// @Tags Багаж
// @Accept json
// @Produce json
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список багажей"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /group/create [post]
func (h *Handler) CreateGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchCode := c.DefaultQuery("searchCode", "")

	var group model.GroupRequest

	if err := c.BindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err := h.UseCase.CreateGroup(authInstance.UserID, group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(searchCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Удаление багажа
// @Description Удаляет багаж по его ID
// @Tags Багаж
// @Produce json
// @Param group_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список багажей"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /group/{group_id}/delete [delete]
func (h *Handler) DeleteGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchCode := c.DefaultQuery("searchCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

	err = h.UseCase.DeleteGroup(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(searchCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Обновление информации о багаже
// @Description Обновляет информацию о багаже по его ID
// @Tags Багаж
// @Accept json
// @Produce json
// @Param group_id path int true "ID багажа"
// @Success 200 {object} model.Group "Информация о багаже"
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id}/update [put]
func (h *Handler) UpdateGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД багажа"}})
		return
	}

	var group model.GroupRequest
	if err := c.BindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не удалось прочитать JSON"})
		return
	}

	err = h.UseCase.UpdateGroup(uint(groupID), authInstance.UserID, group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedGroup, err := h.UseCase.GetGroupByID(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": updatedGroup})
}

// @Summary Добавление багажа к доставке
// @Description Добавляет багаж к доставке по его ID
// @Tags Багаж
// @Produce json
// @Param group_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.GroupsGetResponse  "Список багажей"
// @Failure 400 {object} model.GroupsGetResponse  "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse  "Внутренняя ошибка сервера"
// @Router /group/{group_id}/feedback [post]
func (h *Handler) AddGroupToFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchCode := c.DefaultQuery("searchCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

	err = h.UseCase.AddGroupToFeedback(uint(groupID), authInstance.UserID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(searchCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Удаление багажа из доставки
// @Description Удаляет багаж из доставки по его ID
// @Tags Багаж
// @Produce json
// @Param group_id path int true "ID багажа"
// @Param searchCode query string false "Код багажа" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список багажей"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /groups/{group_id}/feedback [post]
func (h *Handler) RemoveGroupFromFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	searchCode := c.DefaultQuery("searchCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

	err = h.UseCase.RemoveGroupFromFeedback(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(searchCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Добавление изображения к багажу
// @Description Добавляет изображение к багажу по его ID
// @Tags Багаж
// @Accept mpfd
// @Produce json
// @Param group_id path int true "ID багажа"
// @Param image formData file true "Изображение багажа"
// @Success 200 {object} model.Group "Информация о багаже с изображением"
// @Success 200 {object} model.Group
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id}/image [post]
func (h *Handler) AddGroupImage(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД багажа"})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимое изображение"})
		return
	}

	file, err := image.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось открыть изображение"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось прочитать изображение в байтах"})
		return
	}

	contentType := image.Header.Get("Content-Type")

	err = h.UseCase.AddGroupImage(uint(groupID), authInstance.UserID, imageBytes, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	group, err := h.UseCase.GetGroupByID(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}
