package userHandler

import (
	"net/http"
	"redditclone/internal/api/server/session"
	"redditclone/internal/api/server/users"
	"redditclone/utils"
	jsonError "redditclone/utils"
)

type UserRepository interface {
	Get(login string) (*users.User, error)
	Add(users.LoginForm) (*users.User, error)
	CheckPasswords(user *users.User, password string) error
}

type UserHandler struct {
	UserRepo UserRepository
	Session  SessionRepository
}

type SessionRepository interface {
	Check(token string) (*session.UserInfo, error)
	Create(u *users.User) (*session.Session, error)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		jsonError.NewJsonError(w, http.StatusBadRequest, "unknown payload")
		return
	}
	fd := &users.LoginForm{}
	if err := utils.GetDataFromBody(w, r, fd); err != nil {
		return
	}
	foundUser, err := h.UserRepo.Get(fd.Username)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	errPassword := h.UserRepo.CheckPasswords(foundUser, fd.Password)
	if errPassword != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	sess, errCreation := h.Session.Create(foundUser)
	if errCreation != nil {
		jsonError.NewJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	session.CreateAndSendNewToken(sess, w)
}

func (h *UserHandler) Registration(w http.ResponseWriter, r *http.Request) {
	fd := &users.LoginForm{}
	if err := utils.GetDataFromBody(w, r, fd); err != nil {
		return
	}
	userData, err1 := h.UserRepo.Add(*fd)
	if err1 != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err1.Error())
		return
	}
	sess, errCreation := h.Session.Create(userData)
	if errCreation != nil {
		jsonError.NewJsonError(w, http.StatusInternalServerError, errCreation.Error())
		return
	}
	session.CreateAndSendNewToken(sess, w)
}
