package postsRepo

import (
	"context"
	"fmt"
	"math"
	"redditclone/internal/api/server/comments"
	"redditclone/internal/api/server/posts"
	"redditclone/internal/api/server/session"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type PostsRepository struct {
	ctx  context.Context
	coll *mongo.Collection
}

func NewPostsRepo(ctx context.Context, collection *mongo.Collection) *PostsRepository {
	return &PostsRepository{
		ctx:  ctx,
		coll: collection,
	}
}

func (pr *PostsRepository) Get(id string) (*posts.Post, error) {
	var result posts.Post
	err := pr.coll.FindOne(pr.ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("post with id %s was not found", id)
	}
	return &result, nil
}

func (pr *PostsRepository) GetAll() ([]posts.Post, error) {
	var result []posts.Post
	cursor, err := pr.coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error while getting all posts: %s", err)
	}
	if err = cursor.All(context.TODO(), &result); err != nil {
		return nil, fmt.Errorf("error while decoding all posts: %s", err)
	}
	return result, nil
}

func (pr *PostsRepository) Add(body *posts.NewPostForm) (*posts.Post, error) {
	newId := uuid.New().String()
	newPost := &posts.Post{
		Id:               newId,
		Category:         body.Category,
		Created:          time.Now().Format("2006-01-02T15:04:05.000Z"),
		Text:             body.Text,
		URL:              body.Url,
		Title:            body.Title,
		Type:             body.Type,
		Author:           body.Author,
		Comments:         make([]*comments.Comment, 0),
		UpvotePercentage: 100,
		Votes:            make([]*posts.VoteInfo, 0),
	}
	_, err := pr.coll.InsertOne(context.TODO(), newPost)
	if err != nil {
		return nil, fmt.Errorf("error while creating new post")
	}
	return pr.Get(newId)
}

func (pr *PostsRepository) GetAllRelCategories(category string) ([]posts.Post, error) {
	var result []posts.Post
	cursor, err := pr.coll.Find(context.TODO(), bson.M{"category": category})
	if err != nil {
		return nil, fmt.Errorf("error while getting all posts: %s", err)
	}
	if err = cursor.All(pr.ctx, &result); err != nil {
		return nil, fmt.Errorf("error while decoding all posts: %s", err)
	}
	return result, nil
}

func (pr *PostsRepository) GetAllRelUser(author *session.UserInfo) ([]posts.Post, error) {
	var result []posts.Post
	cursor, err := pr.coll.Find(context.TODO(), bson.M{"author": author})
	if err != nil {
		return nil, fmt.Errorf("error while getting all posts: %s", err)
	}
	if err = cursor.All(pr.ctx, &result); err != nil {
		return nil, fmt.Errorf("error while decoding all posts: %s", err)
	}
	return result, nil
}

func (pr *PostsRepository) Upvote(id string, author *session.UserInfo) (*posts.Post, error) {
	post, err := pr.Get(id)
	if err != nil {
		return nil, err
	}
	post.Votes = append(post.Votes, &posts.VoteInfo{User: author.Id, Vote: 1})
	post.UpvotePercentage = countUpvotePercentage(post.Votes)
	_, err = pr.coll.ReplaceOne(context.TODO(), bson.M{"id": id}, post)
	if err != nil {
		return nil, fmt.Errorf("couldn't rewrite post with id %s: %s", id, err.Error())
	}
	return pr.Get(id)
}

func (pr *PostsRepository) DownVote(id string, author *session.UserInfo) (*posts.Post, error) {
	post, err := pr.Get(id)
	if err != nil {
		return nil, err
	}
	post.Votes = append(post.Votes, &posts.VoteInfo{User: author.Id, Vote: -1})
	post.UpvotePercentage = countUpvotePercentage(post.Votes)
	_, err = pr.coll.ReplaceOne(context.TODO(), bson.M{"id": id}, post)
	if err != nil {
		return nil, fmt.Errorf("couldn't downVote post with id %s: %s", id, err.Error())
	}
	return pr.Get(id)
}

func (pr *PostsRepository) Unvote(id string, author *session.UserInfo) (*posts.Post, error) {
	post, err := pr.Get(id)
	if err != nil {
		return nil, err
	}
	for index, val := range post.Votes {
		if val.User == author.Id {
			post.Votes[index] = post.Votes[len(post.Votes)-1]
			post.Votes = post.Votes[:len(post.Votes)-1]
			post.UpvotePercentage = countUpvotePercentage(post.Votes)
		}
	}
	_, err = pr.coll.ReplaceOne(context.TODO(), bson.M{"id": id}, post)
	if err != nil {
		return nil, fmt.Errorf("couldn't unvote post with id %s: %s", id, err.Error())
	}
	return pr.Get(id)
}

func (pr *PostsRepository) Delete(id string) error {
	_, err := pr.coll.DeleteOne(context.TODO(), bson.M{"id": id})
	return err
}

func countUpvotePercentage(arr []*posts.VoteInfo) int {
	allCount := len(arr)
	var upvotedCount int
	for _, item := range arr {
		if item.Vote == 1 {
			upvotedCount++
		}
	}
	if upvotedCount != 0 {
		return int(math.RoundToEven((float64(upvotedCount) / float64(allCount)) * 100))
	} else {
		return 0
	}
}

func (pr *PostsRepository) AddComment(id string, comment *comments.Comment) (*posts.Post, error) {
	post, err := pr.Get(id)
	if err != nil {
		return nil, err
	}
	post.Comments = append(post.Comments, comment)
	_, err = pr.coll.ReplaceOne(context.TODO(), bson.M{"id": id}, post)
	if err != nil {
		return nil, fmt.Errorf("couldn't add comment to post with id %s: %s", id, err.Error())
	}
	return post, nil
}

func (pr *PostsRepository) DeleteComment(id string, messageId string) (*posts.Post, error) {
	post, err := pr.Get(id)
	if err != nil {
		return nil, err
	}
	for inx, comment := range post.Comments {
		if comment.Id == messageId {
			post.Comments = append(post.Comments[:inx], post.Comments[inx+1:]...)
		}
	}
	_, err = pr.coll.ReplaceOne(context.TODO(), bson.M{"id": id}, post)
	if err != nil {
		return nil, fmt.Errorf("couldn't delete comment to post with id %s: %s", id, err.Error())
	}
	return pr.Get(id)
}
