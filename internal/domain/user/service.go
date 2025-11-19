package user

type UserService interface {
	CreateUser(username, password string) (*User, error);
	GetUserByID(id string) (*User, error);
}
