package entity

type Group struct {
	Id string
	Title string
	Creator string // group creator
	CreatedAt int64
}

type GroupUser struct {
	GroupId string
	UserId string
	CreatedAt int64
}