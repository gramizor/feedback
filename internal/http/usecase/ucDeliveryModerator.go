package usecase

import (
	"errors"
	"strings"

	"rest-apishka/internal/model"
)

func (uc *UseCase) GetFeedbacksModerator(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus string, moderatorID uint) ([]model.FeedbackRequest, error) {
	searchFlightNumber = strings.ToUpper(searchFlightNumber + "%")
	feedbackStatus = strings.ToLower(feedbackStatus + "%")

	if moderatorID <= 0 {
		return nil, errors.New("недопустимый ИД модератора")
	}

	feedbacks, err := uc.Repository.GetFeedbacksModerator(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus, moderatorID)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (uc *UseCase) GetFeedbackByIDModerator(feedbackID, moderatorID uint) (model.FeedbackGetResponse, error) {
	if feedbackID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД модератора")
	}

	feedbacks, err := uc.Repository.GetFeedbackByIDModerator(feedbackID, moderatorID)
	if err != nil {
		return model.FeedbackGetResponse{}, err
	}

	return feedbacks, nil
}

func (uc *UseCase) UpdateFlightNumberModerator(feedbackID, moderatorID uint, flightNumber model.FeedbackUpdateFlightNumberRequest) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if len(flightNumber.FlightNumber) != 6 {
		return errors.New("недопустимый номер рейса")
	}

	err := uc.Repository.UpdateFlightNumberModerator(feedbackID, moderatorID, flightNumber)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFeedbackStatusModerator(feedbackID, moderatorID uint, feedbackStatus model.FeedbackUpdateStatusRequest) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД доставки")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if feedbackStatus.FeedbackStatus != model.FEEDBACK_STATUS_COMPLETED && feedbackStatus.FeedbackStatus != model.FEEDBACK_STATUS_REJECTED {
		return errors.New("текущий статус доставки уже завершен или отклонен")
	}

	err := uc.Repository.UpdateFeedbackStatusModerator(feedbackID, moderatorID, feedbackStatus)
	if err != nil {
		return err
	}

	return nil
}
