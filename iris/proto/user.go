package proto

type UserID uint64

type User struct {
	UserID UserID
	Name   string
	Avatar string
}

type UserAuthREQ struct {
	Token string
}
