package entity

type Group struct {
	GroupId string
	Title string
	Creator string // group creator
	CreatedAt int64
}

type GroupUser struct {
	Id int64
	GroupId string
	UserId string
	CreatedAt int64
}