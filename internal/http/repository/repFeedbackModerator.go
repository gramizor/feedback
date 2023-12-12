package repository

import (
	"errors"
	"time"

	"rest-apishka/internal/model"
)

func (r *Repository) GetFeedbacksModerator(startFormationDate, endFormationDate, feedbackStatus string, moderatorID uint) ([]model.FeedbackRequest, error) {
	query := r.db.Table("feedbacks").
		Select("DISTINCT feedbacks.feedback_id, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status LIKE ? AND feedbacks.moderator_id = ? AND feedbacks.feedback_status != ?", feedbackStatus, moderatorID, model.FEEDBACK_STATUS_DELETED)

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
		Select("feedbacks.feedback_id, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status != ? AND feedbacks.feedback_id = ? AND feedbacks.moderator_id = ?", model.FEEDBACK_STATUS_DELETED, feedbackID, moderatorID).
		Scan(&feedback).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения опроса по ИД")
	}

	var groups []model.Group
	if err := r.db.
		Table("groups").
		Joins("JOIN feedback_groups ON groups.group_id = feedback_groups.group_id").
		Where("feedback_groups.feedback_id = ?", feedback.FeedbackID).
		Scan(&groups).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения групп для опроса")
	}

	feedback.Groups = groups

	return feedback, nil
}

func (r *Repository) UpdateFeedbackStatusModerator(feedbackID, moderatorID uint, feedbackStatus model.FeedbackUpdateStatusRequest) error {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").
		Where("feedback_id = ? AND moderator_id = ? AND feedback_status = ?", feedbackID, moderatorID, model.FEEDBACK_STATUS_WORK).
		First(&feedback).
		Error; err != nil {
		return errors.New("опрос не найден или не принадлежит указанному модератору")
	}

	feedback.FeedbackStatus = feedbackStatus.FeedbackStatus
	feedback.CompletionDate = time.Now()

	if err := r.db.Save(&feedback).Error; err != nil {
		return errors.New("ошибка обновления статуса опроса в БД")
	}

	return nil
}
