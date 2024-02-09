package usecase

import (
	"errors"
	"strings"

	"rest-apishka/internal/model"
)

func (uc *UseCase) GetFeedbacksUser(startFormationDate, endFormationDate, feedbackStatus string, userID uint) ([]model.FeedbackRequest, error) {
	feedbackStatus = strings.ToLower(feedbackStatus + "%")

	feedbacks, err := uc.Repository.GetFeedbacksUser(startFormationDate, endFormationDate, feedbackStatus, userID)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (uc *UseCase) GetFeedbackByIDUser(feedbackID, userID uint) (model.FeedbackGetResponse, error) {
	if feedbackID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД опроса")
	}

	feedbacks, err := uc.Repository.GetFeedbackByIDUser(feedbackID, userID)
	if err != nil {
		return model.FeedbackGetResponse{}, err
	}

	return feedbacks, nil
}

func (uc *UseCase) DeleteFeedbackUser(feedbackID, userID uint) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД опроса")
	}

	err := uc.Repository.DeleteFeedbackUser(feedbackID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UseCase) UpdateFeedbackStatusUser(feedbackID, userID uint) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД опроса")
	}

	err := uc.Repository.UpdateFeedbackStatusUser(feedbackID, userID)
	if err != nil {
		return err
	}

	return nil
}
