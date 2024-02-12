package usecase

import (
	"errors"
	"strings"

	"rest-apishka/internal/model"
)

func (uc *UseCase) GetFeedbacksModerator(startFormationDate, endFormationDate, feedbackStatus string) ([]model.FeedbackRequest, error) {
	feedbackStatus = strings.ToLower(feedbackStatus + "%")

	feedbacks, err := uc.Repository.GetFeedbacksModerator(startFormationDate, endFormationDate, feedbackStatus)
	if err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (uc *UseCase) GetFeedbackByIDModerator(feedbackID uint) (model.FeedbackGetResponse, error) {
	if feedbackID <= 0 {
		return model.FeedbackGetResponse{}, errors.New("недопустимый ИД опроса")
	}

	feedbacks, err := uc.Repository.GetFeedbackByIDModerator(feedbackID)
	if err != nil {
		return model.FeedbackGetResponse{}, err
	}

	return feedbacks, nil
}

func (uc *UseCase) UpdateFeedbackStatusModerator(feedbackID, moderatorID uint, feedbackStatus model.FeedbackUpdateStatusRequest) error {
	if feedbackID <= 0 {
		return errors.New("недопустимый ИД опроса")
	}
	if moderatorID <= 0 {
		return errors.New("недопустимый ИД модератора")
	}
	if feedbackStatus.FeedbackStatus != model.FEEDBACK_STATUS_COMPLETED && feedbackStatus.FeedbackStatus != model.FEEDBACK_STATUS_REJECTED {
		return errors.New("текущий статус опроса уже завершен или отклонен")
	}

	err := uc.Repository.UpdateFeedbackStatusModerator(feedbackID, moderatorID, feedbackStatus)
	if err != nil {
		return err
	}

	return nil
}
