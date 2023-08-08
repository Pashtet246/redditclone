package users

type User struct {
	ID       string `json:"id" db:"id"`
	Username string `json:"username" db:"login_id"`
	Password string `json:"password" db:"password"`
}

type LoginForm struct {
	Username string `json:"username" db:"login_id"`
	Password string `json:"password" db:"password"`
}
