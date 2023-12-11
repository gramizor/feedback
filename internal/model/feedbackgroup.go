package model

type FeedbackGroup struct {
	FeedbackID uint `gorm:"type:serial;primaryKey;index" json:"feedback_id"`
	GroupID    uint `gorm:"type:serial;primaryKey;index" json:"group_id"`
}
