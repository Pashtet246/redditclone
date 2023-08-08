package comments

import (
	"redditclone/internal/api/server/session"
)

type Comment struct {
	Author  *session.UserInfo `json:"author" bson:"author"`
	Body    string            `json:"body" bson:"body"`
	Created string            `json:"created" bson:"created"`
	Id      string            `json:"id" bson:"id,omitempty"`
}

type CommentForm struct {
	Comment string `json:"comment"`
}

type CommentRepository interface {
	Get(id string) (*Comment, error)
	Add(body string, author *session.UserInfo) (*Comment, error)
	Delete(id string, author *session.UserInfo) error
}
