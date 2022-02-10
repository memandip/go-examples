package mysql

type Admin struct {
	Id             int64
	Name           string
	Email          string
	Remember_token string
}

type Blog struct {
	Id          int64
	Title       string
	Description string
}
