package user

type Service interface {
	GetAllUsers() ([]User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	ChangeUserPassword(username, newPassword string) error
	DeleteUser(username string) error
	GetAllRoles() ([]Role, error)
	CreateRole(role *Role) error
	UpdateRole(role *Role) error
	DeleteRole(roleID string) error
}

type service struct {
	// Add dependencies here, e.g., a user repository
}

func NewService() Service {
	return &service{}
}

func (s *service) GetAllUsers() ([]User, error) {
	return []User{{Username: "testuser"}}, nil
}

func (s *service) CreateUser(user *User) error {
	return nil
}

func (s *service) UpdateUser(user *User) error {
	return nil
}

func (s *service) ChangeUserPassword(username, newPassword string) error {
	return nil
}

func (s *service) DeleteUser(username string) error {
	return nil
}

func (s *service) GetAllRoles() ([]Role, error) {
	return []Role{{Name: "admin"}}, nil
}

func (s *service) CreateRole(role *Role) error {
	return nil
}

func (s *service) UpdateRole(role *Role) error {
	return nil
}

func (s *service) DeleteRole(roleID string) error {
	return nil
}
