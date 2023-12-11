package repository

import (
	"errors"
	"time"

	"rest-apishka/internal/model"
)

func (r *Repository) GetFeedbacksModerator(searchFlightNumber, startFormationDate, endFormationDate, feedbackStatus string, moderatorID uint) ([]model.FeedbackRequest, error) {
	query := r.db.Table("feedbacks").
		Select("DISTINCT feedbacks.feedback_id, feedbacks.flight_number, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN feedback_groups ON feedbacks.feedback_id = feedback_groups.feedback_id").
		Joins("JOIN groups ON groups.group_id = feedback_groups.group_id").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status LIKE ? AND feedbacks.flight_number LIKE ? AND feedbacks.moderator_id = ? AND feedbacks.feedback_status != ?", feedbackStatus, searchFlightNumber, moderatorID, model.FEEDBACK_STATUS_DELETED)

	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("feedbacks.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	var feedbacks []model.FeedbackRequest
	if err := query.Scan(&feedbacks).Error; err != nil {
		return nil, errors.New("ошибка получения доставок")
	}

	return feedbacks, nil
}

func (r *Repository) GetFeedbackByIDModerator(feedbackID, moderatorID uint) (model.FeedbackGetResponse, error) {
	var feedback model.FeedbackGetResponse

	if err := r.db.
		Table("feedbacks").
		Select("feedbacks.feedback_id, feedbacks.flight_number, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status != ? AND feedbacks.feedback_id = ? AND feedbacks.moderator_id = ?", model.FEEDBACK_STATUS_DELETED, feedbackID, moderatorID).
		Scan(&feedback).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения доставки по ИД")
	}

	var groups []model.Group
	if err := r.db.
		Table("groups").
		Joins("JOIN feedbacks_groups ON groups.group_id = feedbacks_groups.group_id").
		Where("feedbacks_groups.feedback_id = ?", feedback.FeedbackID).
		Scan(&groups).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения багажей для доставки")
	}

	feedback.Groups = groups

	return feedback, nil
}

func (r *Repository) UpdateFlightNumberModerator(feedbackID uint, moderatorID uint, flightNumber model.FeedbackUpdateFlightNumberRequest) error {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").
		Where("feedback_id = ? AND moderator_id = ?", feedbackID, moderatorID).
		First(&feedback).
		Error; err != nil {
		return errors.New("доставка не найдена или не принадлежит указанному модератору")
	}

	if err := r.db.Table("feedbacks").
		Model(&feedback).
		Update("flight_number", flightNumber.FlightNumber).
		Error; err != nil {
		return errors.New("ошибка обновления номера рейса")
	}

	return nil
}

func (r *Repository) UpdateFeedbackStatusModerator(feedbackID, moderatorID uint, feedbackStatus model.FeedbackUpdateStatusRequest) error {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").
		Where("feedback_id = ? AND moderator_id = ? AND feedback_status = ?", feedbackID, moderatorID, model.FEEDBACK_STATUS_WORK).
		First(&feedback).
		Error; err != nil {
		return errors.New("доставка не найдена или не принадлежит указанному модератору")
	}

	feedback.FeedbackStatus = feedbackStatus.FeedbackStatus
	feedback.CompletionDate = time.Now()

	if err := r.db.Save(&feedback).Error; err != nil {
		return errors.New("ошибка обновления статуса доставки в БД")
	}

	return nil
}
