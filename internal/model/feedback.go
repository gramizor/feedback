package model

import "time"

type Feedback struct {
	FeedbackID     uint      `gorm:"type:serial;primarykey" json:"feedback_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	FeedbackStatus string    `json:"feedback_status"`
	UserID         uint      `json:"user_id"`
	ModeratorID    uint      `json:"moderator_id"`
}

type FeedbackRequest struct {
	FeedbackID     uint      `json:"feedback_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	FeedbackStatus string    `json:"feedback_status"`
	FullName       string    `json:"full_name"`
}

type FeedbackGetResponse struct {
	FeedbackID     uint      `json:"feedback_id"`
	FlightNumber   string    `json:"flight_number"`
	CreationDate   time.Time `json:"creation_date"`
	FormationDate  time.Time `json:"formation_date"`
	CompletionDate time.Time `json:"completion_date"`
	FeedbackStatus string    `json:"feedback_status"`
	FullName       string    `json:"full_name"`
	Groups         []Group   `json:"groups"`
}

type FeedbackUpdateFlightNumberRequest struct {
	FlightNumber string `json:"flight_number"`
}

type FeedbackUpdateStatusRequest struct {
	FeedbackStatus string `json:"feedback_status"`
}
