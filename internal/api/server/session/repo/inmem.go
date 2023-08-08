package sessionRepo

import (
	"encoding/json"
	"fmt"
	"redditclone/internal/api/server/session"
	"redditclone/internal/api/server/users"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var SessionDuration int64 = 600

type SessionRepository struct {
	redisClient *redis.Client
}

func NewSessionRepo(client *redis.Client) *SessionRepository {
	return &SessionRepository{
		redisClient: client,
	}
}

func (sr *SessionRepository) Check(token string) (*session.UserInfo, error) {
	payload, err := session.ParseJWTToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while parsing JWT")
	}
	return &session.UserInfo{Id: payload.User.Id, Username: payload.User.Username}, nil
}

func (sr *SessionRepository) Create(u *users.User) (*session.Session, error) {
	createdTime := time.Now().Unix()
	userInfo := session.UserInfo{
		Username: u.Username,
		Id:       u.ID,
	}
	newId := uuid.New().String()
	expiration := time.Duration(600 * time.Second)
	jsonData, err := json.Marshal(&session.Session{
		Iat:  createdTime,
		Exp:  createdTime + SessionDuration,
		User: &userInfo,
		Id:   newId,
	})
	if err != nil {
		return nil, fmt.Errorf("error while marshal jsonData: %s", err.Error())
	}
	_, err = sr.redisClient.Set(newId, jsonData, expiration).Result()
	if err != nil {
		return nil, fmt.Errorf("error while redis set jsonData: %s", err.Error())
	}

	sessionInfo := session.Session{}

	err = sr.redisClient.Get(newId).Scan(&sessionInfo)
	if err != nil {
		return nil, fmt.Errorf("error while redis get jsonData: %s", err.Error())
	}
	return &sessionInfo, nil
}
