package delivery

import (
	"io"
	"net/http"
	"strconv"

	"rest-apishka/internal/auth"
	"rest-apishka/internal/model"

	"github.com/gin-gonic/gin"
)

// @Summary Получение списка групп
// @Description Возращает список всех активных групп
// @Tags Группа
// @Produce json
// @Param groupCode query string false "Код группы" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список групп"
// @Failure 500 {object} model.GroupsGetResponse "Ошибка сервера"
// @Router /group [get]
func (h *Handler) GetGroups(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	groups, err := h.UseCase.GetGroupsPaged(groupCode, authInstance.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Получение списка групп с пагинацией
// @Description Возвращает список всех активных групп с использованием пагинации
// @Tags Группа
// @Produce json
// @Param groupCode query string false "Код группы" Format(email)
// @Param page query int false "Номер страницы" Format(email)
// @Param pageSize query int false "Размер страницы" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список групп"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /group/paginate [get]
func (h *Handler) GetGroupsPaged(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "10")

	groups, err := h.UseCase.GetGroupsPaged(groupCode, authInstance.UserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Получение группы по ID
// @Description Возвращает информацию о группе по его ID
// @Tags Группа
// @Produce json
// @Param group_id path int true "ID группы"
// @Success 200 {object} model.Group "Информация о группе"
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id} [get]
func (h *Handler) GetGroupByID(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД группы"})
		return
	}

	group, err := h.UseCase.GetGroupByID(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"group": group})
}

// @Summary Создание новой группы
// @Description Создает новую группу с предоставленными данными
// @Tags Группа
// @Accept json
// @Produce json
// @Param groupCode query string false "Код группы" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список групп"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /group/create [post]
func (h *Handler) CreateGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")

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

	groups, err := h.UseCase.GetGroups(groupCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Удаление группы
// @Description Удаляет группу по ее ID
// @Tags Группа
// @Produce json
// @Param group_id path int true "ID группы"
// @Param groupCode query string false "Код группы" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список групп"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /group/{group_id}/delete [delete]
func (h *Handler) DeleteGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД группы"})
		return
	}

	err = h.UseCase.DeleteGroup(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(groupCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Обновление информации о группе
// @Description Обновляет информацию о группе по ее ID
// @Tags Группа
// @Accept json
// @Produce json
// @Param group_id path int true "ID группы"
// @Success 200 {object} model.Group "Информация о группе"
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id}/update [put]
func (h *Handler) UpdateGroup(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"error": "недопустимый ИД группы"}})
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

// @Summary Добавление группы к опросу
// @Description Добавляет группу к опросу по ее ID
// @Tags Группа
// @Produce json
// @Param group_id path int true "ID группы"
// @Param groupCode query string false "Код группы" Format(email)
// @Success 200 {object} model.GroupsGetResponse  "Список групп"
// @Failure 400 {object} model.GroupsGetResponse  "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse  "Внутренняя ошибка сервера"
// @Router /group/{group_id}/feedback [post]
func (h *Handler) AddGroupToFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД группы"})
		return
	}

	err = h.UseCase.AddGroupToFeedback(uint(groupID), authInstance.UserID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(groupCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Удаление группы из опроса
// @Description Удаляет группу из опроса по ее ID
// @Tags Группа
// @Produce json
// @Param group_id path int true "ID группы"
// @Param groupCode query string false "Код группы" Format(email)
// @Success 200 {object} model.GroupsGetResponse "Список групп"
// @Failure 400 {object} model.GroupsGetResponse "Некорректный запрос"
// @Failure 500 {object} model.GroupsGetResponse "Внутренняя ошибка сервера"
// @Router /groups/{group_id}/feedback [post]
func (h *Handler) RemoveGroupFromFeedback(c *gin.Context) {
	authInstance := auth.GetAuthInstance()
	groupCode := c.DefaultQuery("groupCode", "")

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД группы"})
		return
	}

	err = h.UseCase.RemoveGroupFromFeedback(uint(groupID), authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	groups, err := h.UseCase.GetGroups(groupCode, authInstance.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}

// @Summary Добавление изображения к группе
// @Description Добавляет изображение к группе по его ID
// @Tags Группа
// @Accept mpfd
// @Produce json
// @Param group_id path int true "ID группы"
// @Param image formData file true "Изображение группы"
// @Success 200 {object} model.Group "Информация о группе с изображением"
// @Success 200 {object} model.Group
// @Failure 400 {object} model.Group "Некорректный запрос"
// @Failure 500 {object} model.Group "Внутренняя ошибка сервера"
// @Router /group/{group_id}/image [post]
func (h *Handler) AddGroupImage(c *gin.Context) {
	authInstance := auth.GetAuthInstance()

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимый ИД группы"})
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
