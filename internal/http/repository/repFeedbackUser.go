package repository

import (
	"errors"
	"time"

	"rest-apishka/internal/model"
)

func (r *Repository) GetFeedbacksUser(startFormationDate, endFormationDate, feedbackStatus string, userID uint) ([]model.FeedbackRequest, error) {
	query := r.db.Table("feedbacks").
		Select("DISTINCT feedbacks.feedback_id, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status LIKE ? AND feedbacks.user_id = ? AND feedbacks.feedback_status != ?", feedbackStatus, userID, model.FEEDBACK_STATUS_DELETED)

	if startFormationDate != "" && endFormationDate != "" {
		query = query.Where("feedbacks.formation_date BETWEEN ? AND ?", startFormationDate, endFormationDate)
	}

	var feedbacks []model.FeedbackRequest
	if err := query.Scan(&feedbacks).Error; err != nil {
		return nil, errors.New("ошибка получения опросов")
	}

	return feedbacks, nil
}

func (r *Repository) GetFeedbackByIDUser(feedbackID, userID uint) (model.FeedbackGetResponse, error) {
	var feedback model.FeedbackGetResponse

	if err := r.db.
		Table("feedbacks").
		Select("feedbacks.feedback_id, feedbacks.creation_date, feedbacks.formation_date, feedbacks.completion_date, feedbacks.feedback_status, users.full_name").
		Joins("JOIN users ON users.user_id = feedbacks.user_id").
		Where("feedbacks.feedback_status != ? AND feedbacks.feedback_id = ? AND feedbacks.user_id = ?", model.FEEDBACK_STATUS_DELETED, feedbackID, userID).
		Scan(&feedback).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения опроса по ИД")
	}

	var groups []model.Group
	if err := r.db.
		Table("groups").
		Joins("JOIN feedback_groups ON groups.group_id = feedback_groups.group_id").
		Where("feedback_groups.feedback_id = ?", feedback.FeedbackID).
		Scan(&groups).Error; err != nil {
		return model.FeedbackGetResponse{}, errors.New("ошибка получения группы для опроса")
	}

	feedback.Groups = groups

	return feedback, nil
}

func (r *Repository) DeleteFeedbackUser(feedbackID, userID uint) error {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").
		Where("feedback_id = ? AND user_id = ?", feedbackID, userID).
		First(&feedback).
		Error; err != nil {
		return errors.New("опрос не найден или не принадлежит указанному преподавателю")
	}

	tx := r.db.Begin()
	if err := tx.Where("feedback_id = ?", feedbackID).Delete(&model.FeedbackGroup{}).Error; err != nil {
		tx.Rollback()
		return errors.New("ошибка удаления связи таблицы feedback_groups")
	}

	err := r.db.Model(&model.Feedback{}).Where("feedback_id = ?", feedbackID).Update("feedback_status", model.FEEDBACK_STATUS_DELETED).Error
	if err != nil {
		return errors.New("ошибка обновления статуса на удален")
	}
	tx.Commit()

	return nil
}

func (r *Repository) UpdateFeedbackStatusUser(feedbackID, userID uint) error {
	var feedback model.Feedback
	if err := r.db.Table("feedbacks").
		Where("feedback_id = ? AND user_id = ? AND feedback_status = ?", feedbackID, userID, model.FEEDBACK_STATUS_DRAFT).
		First(&feedback).
		Error; err != nil {
		return errors.New("опрос не найден, или не принадлежит указанному преподавателю, или не имеет статус черновика")
	}

	currentTime := time.Now()
	feedback.FeedbackStatus = model.FEEDBACK_STATUS_WORK
	feedback.FormationDate = &currentTime

	if err := r.db.Save(&feedback).Error; err != nil {
		return errors.New("ошибка обновления статуса опроса в БД")
	}

	return nil
}
