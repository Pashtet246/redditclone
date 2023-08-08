package repo

import (
	"database/sql"
	"fmt"
	"redditclone/internal/api/server/users"

	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (ur *UsersRepository) Get(login string) (*users.User, error) {
	var user users.User
	err := ur.db.Get(&user, "SELECT * FROM users WHERE login_id=?", login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user were found")
		}
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	return &user, nil
}

func (ur *UsersRepository) Add(u users.LoginForm) (*users.User, error) {
	ur.db.MustExec("INSERT INTO users (login_id, password) VALUES (?, ?)", u.Username, u.Password)
	return ur.Get(u.Username)
}

func (ur *UsersRepository) CheckPasswords(user *users.User, password string) error {
	if user.Password != password {
		return fmt.Errorf("password is incorrect")
	}
	return nil
}
