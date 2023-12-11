// Package model ...
package model

// Group представляет информацию о багаже.
type Group struct {
	GroupID        uint    `gorm:"type:serial;primarykey" json:"group_id"`
	GroupCode      string  `json:"group_code" example:"ABC123"`
	Weight         float32 `json:"weight" example:"23.5"`
	Size           string  `json:"size" example:"large"`
	GroupStatus    string  `json:"group_status" example:"checked"`
	GroupType      string  `json:"group_type" example:"suitcase"`
	OwnerName      string  `json:"owner_name" example:"John Doe"`
	PasportDetails string  `json:"pasport_details" example:"123456789"`
	Airline        string  `json:"airline" example:"AirlineX"`
	PhotoURL       string  `json:"photo" example:"http://example.com/group.jpg"`
}

// GroupRequest представляет запрос на создание багажа.
type GroupRequest struct {
	GroupCode      string  `json:"group_code" example:"ABC123"`
	Weight         float32 `json:"weight" example:"23.5"`
	Size           string  `json:"size" example:"large"`
	GroupType      string  `json:"group_type" example:"suitcase"`
	OwnerName      string  `json:"owner_name" example:"John Doe"`
	PasportDetails string  `json:"pasport_details" example:"123456789"`
	Airline        string  `json:"airline" example:"AirlineX"`
}

// GroupsGetResponse представляет ответ с информацией о багажах и идентификаторе доставки.
type GroupsGetResponse struct {
	Groups     []Group `json:"groups"`
	FeedbackID uint    `json:"feedback_id" example:"1"`
}
