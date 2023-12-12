// Package model ...
package model

// Group представляет информацию о группе.
type Group struct {
	GroupID     uint   `gorm:"type:serial;primarykey" json:"group_id"`
	GroupCode   string `json:"group_code" example:"RT5-51B"`
	Contacts    string `json:"contacts" example:"+7(999)999-99-99"`
	Course      int    `json:"course" example:"1"`
	Students    int    `json:"students" example:"23"`
	GroupStatus string `json:"group_status" example:"обучается"`
	Photo       string `json:"photo" example:"http://example.com/group.jpg"`
}

// GroupRequest представляет запрос на создание группы.
type GroupRequest struct {
	GroupCode string `json:"group_code" example:"RT5-51B"`
	Contacts  string `json:"contacts" example:"+7(999)999-99-99"`
	Course    int    `json:"course" example:"1"`
	Students  int    `json:"students" example:"23"`
}

// GroupsGetResponse представляет ответ с информацией о группах и идентификаторе опроса.
type GroupsGetResponse struct {
	Groups     []Group `json:"groups"`
	FeedbackID uint    `json:"feedback_id" example:"1"`
}
