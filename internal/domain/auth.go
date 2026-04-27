package domain

type User struct {
	ID    uint   `json:"id"`
	Phone string `json:"phone"`
	Name  string `json:"name"`
	Role  string `json:"role"` // guest, user, provider, admin
}

type AuthRepository interface {
	// TODO: define methods
}

type AuthService interface {
	// TODO: define methods
}
