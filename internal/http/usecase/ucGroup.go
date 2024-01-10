package usecase

import (
	"errors"
	"strconv"
	"strings"

	"rest-apishka/internal/model"
)

type GroupUseCase interface {
}

func (uc *UseCase) GetGroups(groupCode string, userID uint) (model.GroupsGetResponse, error) {
	if userID <= 0 {
		return model.GroupsGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	groupCode = strings.ToUpper(groupCode + "%")

	groups, err := uc.Repository.GetGroups(groupCode, userID)
	if err != nil {
		return model.GroupsGetResponse{}, err
	}

	return groups, nil
}

func (uc *UseCase) GetGroupsPaged(groupCode string, courseNumber int, userID uint, page string, pageSize string) (model.GroupsGetResponse, error) {
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return model.GroupsGetResponse{}, errors.New("некорректное значение страницы")
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		return model.GroupsGetResponse{}, errors.New("некорректное значение размера страницы")
	}

	if userID <= 0 {
		return model.GroupsGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	groupCode = strings.ToUpper(groupCode + "%")

	groups, err := uc.Repository.GetGroupsPaged(groupCode, courseNumber, userID, pageNum, pageSizeNum)
	if err != nil {
		return model.GroupsGetResponse{}, err
	}

	groupResponse := model.GroupsGetResponse{
		Groups:     groups.Groups,
		FeedbackID: groups.FeedbackID,
	}

	return groupResponse, nil
}

func (uc *UseCase) GetGroupByID(groupID, userID uint) (model.Group, error) {
	if groupID <= 0 {
		return model.Group{}, errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return model.Group{}, errors.New("недопустимый ИД пользователя")
	}

	group, err := uc.Repository.GetGroupByID(groupID, userID)
	if err != nil {
		return model.Group{}, err
	}

	return group, nil
}

func (uc *UseCase) CreateGroup(userID uint, requestGroup model.GroupRequest) error {
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if requestGroup.GroupCode == "" {
		return errors.New("название группы должно быть заполнено")
	}
	if requestGroup.Contacts == "" {
		return errors.New("контакты группы должны быть заполнены")
	}
	if requestGroup.Course == 0 {
		return errors.New("номер курса группы должен быть заполнен")
	}
	if requestGroup.Students == 0 {
		return errors.New("количество студентов должно быть заполнено")
	}

	group := model.Group{
		GroupCode:   requestGroup.GroupCode,
		Contacts:    requestGroup.Contacts,
		Course:      requestGroup.Course,
		Students:    requestGroup.Students,
		GroupStatus: model.GROUP_STATUS_ACTIVE,
	}

	err := uc.Repository.CreateGroup(userID, group)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DeleteGroup(groupID, userID uint) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeleteGroup(groupID, userID)
	if err != nil {
		return err
	}

	err = uc.Repository.RemoveServiceImage(groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateGroup(groupID, userID uint, requestGroup model.GroupRequest) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	group := model.Group{
		GroupCode: requestGroup.GroupCode,
		Contacts:  requestGroup.Contacts,
		Course:    requestGroup.Course,
		Students:  requestGroup.Students,
	}

	err := uc.Repository.UpdateGroup(groupID, userID, group)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddGroupToFeedback(groupID, userID, moderatorID uint) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}

	err := uc.Repository.AddGroupToFeedback(groupID, userID, moderatorID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) RemoveGroupFromFeedback(groupID, userID uint) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.RemoveGroupFromFeedback(groupID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddGroupImage(groupID, userID uint, imageBytes []byte, ContentType string) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД группы")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if imageBytes == nil {
		return errors.New("недопустимый imageBytes изображения")
	}
	if ContentType == "" {
		return errors.New("недопустимый ContentType изображения")
	}

	imageURL, err := uc.Repository.UploadServiceImage(groupID, userID, imageBytes, ContentType)
	if err != nil {
		return err
	}

	err = uc.Repository.AddGroupImage(groupID, userID, imageURL)
	if err != nil {
		return err
	}

	return nil
}
