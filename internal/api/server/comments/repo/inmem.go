package commentsRepo

import (
	"context"
	"fmt"
	"redditclone/internal/api/server/comments"
	"redditclone/internal/api/server/session"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type CommentsRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewCommentsRepo(ctx context.Context, collection *mongo.Collection) *CommentsRepository {
	return &CommentsRepository{
		ctx:  ctx,
		coll: collection,
	}
}

func (cr *CommentsRepository) Get(id string) (*comments.Comment, error) {
	var result comments.Comment
	err := cr.coll.FindOne(cr.ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("comment with id %s was not found", id)
	}
	return &result, nil
}

func (cr *CommentsRepository) Add(body string, author *session.UserInfo) (*comments.Comment, error) {
	newId := uuid.New().String()
	newComment := &comments.Comment{
		Body:    body,
		Created: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Author:  author,
		Id:      newId,
	}
	_, err := cr.coll.InsertOne(context.TODO(), newComment)
	if err != nil {
		return nil, fmt.Errorf("error while creating new comment")
	}
	return cr.Get(newId)
}

func (cr *CommentsRepository) Delete(id string, author *session.UserInfo) error {
	comment, err := cr.Get(id)
	if err != nil {
		return err
	}
	if comment.Author.Id != author.Id {
		return fmt.Errorf("comment with id %s author is wrong", id)
	}
	_, err = cr.coll.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}
