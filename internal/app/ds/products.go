package ds

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	Id       uint `gorm:"primarykey"`
	Name     string
	Contacts string
	Course   int
	Image    string
	Status   string
	Href     string
	Students int
	Requests []Request `gorm:"many2many:group_request"`
}

type Request struct {
	Id             uint `gorm:"primarykey"`
	Status         string
	CreationDate   time.Time
	FormationDate  time.Time
	CompletionDate time.Time
	CreatorID      uint  `gorm:"index"`
	ModeratorID    uint  `gorm:"index"`
	GroupID        uint  `gorm:"index"`
	Creator        User  `gorm:"foreignKey:CreatorID"`
	Moderator      User  `gorm:"foreignKey:ModeratorID"`
	Group          Group `gorm:"foreignKey:GroupID"`
}

type User struct {
	Id        uint `gorm:"primarykey"`
	Username  string
	Password  string
	Role      string
	Requests  []Request `gorm:"foreignKey:CreatorID;references:Id"`
	Moderates []Request `gorm:"foreignKey:ModeratorID;references:Id"`
}

func (group *Group) BeforeCreate(tx *gorm.DB) (err error) {
	if group.Status == "" {
		group.Status = "active" // Задайте значение по умолчанию, например, "active"
	}
	return nil
}
