package usecase

import (
	"errors"
	"strings"

	"rest-apishka/internal/model"
)

type GroupUseCase interface {
}

func (uc *UseCase) GetGroups(searchCode string, userID uint) (model.GroupsGetResponse, error) {
	if userID <= 0 {
		return model.GroupsGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	searchCode = strings.ToUpper(searchCode + "%")

	groups, err := uc.Repository.GetGroups(searchCode, userID)
	if err != nil {
		return model.GroupsGetResponse{}, err
	}

	return groups, nil
}

func (uc *UseCase) GetGroupByID(groupID, userID uint) (model.Group, error) {
	if groupID <= 0 {
		return model.Group{}, errors.New("недопустимый ИД багажа")
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
		return errors.New("код багажа должен быть заполнен")
	}
	if requestGroup.Weight == 0 {
		return errors.New("вес багажа должен быть заполнен")
	}
	if requestGroup.Size == "" {
		return errors.New("размер багажа должен быть заполнен")
	}
	if requestGroup.GroupType == "" {
		return errors.New("тип багажа должен быть заполнен")
	}
	if requestGroup.OwnerName == "" {
		return errors.New("владелец багажа должен быть заполнен")
	}
	if requestGroup.PasportDetails == "" {
		return errors.New("паспортные данные владельца багажа должны быть заполнен")
	}
	if requestGroup.Airline == "" {
		return errors.New("авиакомпания должна быть заполнена")
	}

	group := model.Group{
		GroupCode:      requestGroup.GroupCode,
		Weight:         requestGroup.Weight,
		Size:           requestGroup.Size,
		GroupType:      requestGroup.GroupType,
		OwnerName:      requestGroup.OwnerName,
		PasportDetails: requestGroup.PasportDetails,
		Airline:        requestGroup.Airline,
		GroupStatus:    model.GROUP_STATUS_ACTIVE,
	}

	err := uc.Repository.CreateGroup(userID, group)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) DeleteGroup(groupID, userID uint) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД багажа")
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
		return errors.New("недопустимый ИД багажа")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	group := model.Group{
		GroupCode:      requestGroup.GroupCode,
		Weight:         requestGroup.Weight,
		Size:           requestGroup.Size,
		GroupType:      requestGroup.GroupType,
		OwnerName:      requestGroup.OwnerName,
		PasportDetails: requestGroup.PasportDetails,
		Airline:        requestGroup.Airline,
	}

	err := uc.Repository.UpdateGroup(groupID, userID, group)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) AddGroupToFeedback(groupID, userID, moderatorID uint) error {
	if groupID <= 0 {
		return errors.New("недопустимый ИД багажа")
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
		return errors.New("недопустимый ИД багажа")
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
		return errors.New("недопустимый ИД багажа")
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
