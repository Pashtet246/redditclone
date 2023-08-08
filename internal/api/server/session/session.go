package session

import (
	"context"
	"encoding/json"
	"errors"
)

type Session struct {
	Iat  int64     `json:"iat"`
	Exp  int64     `json:"exp"`
	User *UserInfo `json:"user"`
	Id   string
}

type UserInfo struct {
	Username string `json:"username"`
	Id       string `json:"id"`
}

var (
	ErrNoAuth = errors.New("no session found")
)

type sessKey string

var SessionKey sessKey = "sessionKey"

func SessionFromContext(ctx context.Context) (*UserInfo, error) {
	sess, ok := ctx.Value(SessionKey).(*UserInfo)
	if !ok || sess == nil {
		return nil, ErrNoAuth
	}
	return sess, nil
}

func (s *Session) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

func (s *Session) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
