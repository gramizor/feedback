package repository

import (
	"errors"
	"time"

	"rest-apishka/internal/model"
)

type GroupRepository interface {
	GetGroups(searchCode string, userID uint) (model.GroupsGetResponse, error)
}

func (r *Repository) GetGroups(searchCode string, userID uint) (model.GroupsGetResponse, error) {
	var feedbackID uint
	if err := r.db.
		Table("feedbacks").
		Select("feedbacks.feedback_id").
		Where("user_id = ? AND feedback_status = ?", userID, model.FEEDBACK_STATUS_DRAFT).
		Take(&feedbackID).Error; err != nil {
	}

	var groups []model.Group
	if err := r.db.Table("groups").
		Where("groups.group_status = ? AND groups.group_code LIKE ?", model.GROUP_STATUS_ACTIVE, searchCode).
		Scan(&groups).Error; err != nil {
		return model.GroupsGetResponse{}, errors.New("ошибка нахождения списка багажа")
	}

	groupResponse := model.GroupsGetResponse{
		Groups:     groups,
		FeedbackID: feedbackID,
	}

	return groupResponse, nil
}

func (r *Repository) GetGroupByID(groupID, userID uint) (model.Group, error) {
	var group model.Group

	if err := r.db.Table("groups").
		Where("group_status = ? AND group_id = ?", model.GROUP_STATUS_ACTIVE, groupID).
		First(&group).Error; err != nil {
		return model.Group{}, errors.New("ошибка при получении активного багажа из БД")
	}

	return group, nil
}

func (r *Repository) CreateGroup(userID uint, group model.Group) error {
	if err := r.db.Create(group).Error; err != nil {
		return errors.New("ошибка создания багажа")
	}

	return nil
}

func (r *Repository) DeleteGroup(groupID, userID uint) error {
	var group model.Group

	if err := r.db.Table("groups").Where("group_id = ? AND group_status = ?", groupID, model.GROUP_STATUS_ACTIVE).First(group).Error; err != nil {
		return errors.New("багаж не найден или уже удален")
	}

	group.GroupStatus = model.GROUP_STATUS_DELETED

	if err := r.db.Table("groups").Save(group).Error; err != nil {
		return errors.New("ошибка при обновлении статуса багажа в БД")
	}
	return nil
}

func (r *Repository) UpdateGroup(groupID, userID uint, group model.Group) error {
	if err := r.db.Table("groups").
		Model(&model.Group{}).
		Where("group_id = ? AND group_status = ?", groupID, model.GROUP_STATUS_ACTIVE).
		Updates(group).Error; err != nil {
		return errors.New("ошибка при обновлении информации о питомце в БД")
	}

	return nil
}

func (r *Repository) AddGroupToFeedback(groupID, userID, moderatorID uint) error {
	var group model.Group

	if err := r.db.Table("groups").
		Where("group_id = ? AND group_status = ?", groupID, model.GROUP_STATUS_ACTIVE).
		First(group).Error; err != nil {
		return errors.New("багаж не найден или удален")
	}

	var feedback model.Feedback

	if err := r.db.Table("feedback").
		Where("feedback_status = ? AND user_id = ?", model.FEEDBACK_STATUS_DRAFT, userID).
		Last(feedback).Error; err != nil {
		feedback = model.Feedback{
			FeedbackStatus: model.FEEDBACK_STATUS_DRAFT,
			CreationDate:   time.Now(),
			UserID:         userID,
			ModeratorID:    moderatorID,
		}

		if err := r.db.Table("feedback").
			Create(feedback).Error; err != nil {
			return errors.New("ошибка создания доставки со статусом черновик")
		}
	}

	feedbackGroup := model.FeedbackGroup{
		GroupID:    groupID,
		FeedbackID: feedback.FeedbackID,
	}

	if err := r.db.Table("feedbacks_groups").
		Create(feedbackGroup).Error; err != nil {
		return errors.New("ошибка при создании связи между доставкой и багажом")
	}

	return nil
}

func (r *Repository) RemoveGroupFromFeedback(groupID, userID uint) error {
	var feedbackGroup model.FeedbackGroup

	if err := r.db.Joins("JOIN feedbacks ON feedback_groups.feedback_id = feedbacks.feedback_id").
		Where("feedback_groups.group_id = ? AND feedbacks.user_id = ? AND feedbacks.feedback_status = ?", groupID, userID, model.FEEDBACK_STATUS_DRAFT).
		First(feedbackGroup).Error; err != nil {
		return errors.New("багаж не принадлежит пользователю или находится не в статусе черновик")
	}

	if err := r.db.Table("feedbacks_groups").
		Delete(feedbackGroup).Error; err != nil {
		return errors.New("ошибка удаления связи между доставкой и багажом")
	}

	return nil
}

func (r *Repository) AddGroupImage(groupID, userID uint, imageURL string) error {
	err := r.db.Table("groups").Where("group_id = ?", groupID).Update("photo", imageURL).Error
	if err != nil {
		return errors.New("ошибка обновления url изображения в БД")
	}

	return nil
}
