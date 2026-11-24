package entities

// User contains user data shared by application services and persistence adapters.
type User struct {
	ID       string
	Name     string
	Username string
	Phone    string
	Password string
}
