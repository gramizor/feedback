package usecase

import (
	"errors"
	"strings"

	"rest-apishka/internal/model"
)

func (uc *UseCase) GetFeedbacksUser(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus string, userID uint) ([]model.FeedbackRequest, error) {
	searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
	feedbackStatus = strings.ToLower(feedbackStatus + "%")

	if userID <= 0 {
		return nil, errors.New("недопустимый ИД пользователя")
	}

	feedbacks, err := uc.Repository.GetFeedbacksUser(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus, userID)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (uc *UseCase) GetFeedbackByIDUser(feedbackID, userID uint) (model.FeedbackGetResponse, error) {
	if feedbackID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД пользователя")
	}

	feedbacks, err := uc.Repository.GetFeedbackByIDUser(feedbackID, userID)
	if err != nil {
		return model.FeedbackGetResponse{}, err
	}

	return feedbacks, nil
}

func (uc *UseCase) DeleteFeedbackUser(feedbackID, userID uint) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.DeleteFeedbackUser(feedbackID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFlightNumberUser(feedbackID, userID uint, flightNumber model.FeedbackUpdateFlightNumberRequest) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}
	if len(flightNumber.FlightNumber) != 6 {
		return errors.New("недопустимый номер рейса")
	}

	err := uc.Repository.UpdateFlightNumberUser(feedbackID, userID, flightNumber)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFeedbackStatusUser(feedbackID, userID uint) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if userID <= 0 {
		return errors.New("недопустимый ИД пользователя")
	}

	err := uc.Repository.UpdateFeedbackStatusUser(feedbackID, userID)
	if err != nil {
		return err
	}

	return nil
}
