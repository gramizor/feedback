package api

type Group struct {
	Id       int
	Name     string
	Contacts string
	Course   int
	Subjects string
	Src      string
	Href     string
	Students []Student
}

type Student struct {
	Name string
	Post string
}
