package session

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jsonError "redditclone/utils"

	"github.com/dgrijalva/jwt-go"
)

var SecretToken = []byte("secretToken")

type LocaleSession struct {
	Iat       int64    `json:"iat"`
	Exp       int64    `json:"exp"`
	SessionId string   `json:"sessionId"`
	User      UserInfo `json:"user"`
}

func (payload LocaleSession) Valid() error {
	if payload.Exp < time.Now().Unix() {
		return fmt.Errorf("token was expired")
	}
	return nil
}

func ParseJWTToken(token string) (*LocaleSession, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return SecretToken, nil
	}
	sess := &LocaleSession{}
	parsedToken, err := jwt.ParseWithClaims(token, sess, hashSecretGetter)
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("bad token")
	}
	return sess, nil
}

func CreateAndSendNewToken(sessionData *Session, w http.ResponseWriter) {
	sess := &LocaleSession{
		Iat:       sessionData.Iat,
		Exp:       sessionData.Exp,
		SessionId: sessionData.Id,
		User:      *sessionData.User,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sess)
	tokenString, err := token.SignedString(SecretToken)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response, err := json.Marshal(map[string]interface{}{
		"token": tokenString,
	})
	if err != nil {
		jsonError.NewJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Write(response)
}
