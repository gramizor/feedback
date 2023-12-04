package ds

import "time"

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
	Creator        User  `gorm:"foreignKey:CreatorID"`
	Moderator      User  `gorm:"foreignKey:ModeratorID"`
	GroupID        uint  `gorm:"index"`
	Groups         Group `gorm:"many2many:group_request"`
}

type User struct {
	Id        uint `gorm:"primarykey"`
	Username  string
	Password  string
	Role      string
	Requests  []Request `gorm:"foreignKey:CreatorID;references:Id"`
	Moderates []Request `gorm:"foreignKey:ModeratorID;references:Id"`
}
