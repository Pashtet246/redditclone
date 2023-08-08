package posts

import (
	"net/http"
	"redditclone/internal/api/server/comments"
	"redditclone/internal/api/server/session"
	userHandler "redditclone/internal/api/server/users/delivery"
	"redditclone/utils"
	jsonError "redditclone/utils"

	"github.com/gorilla/mux"
)

type Post struct {
	Author           *session.UserInfo   `json:"author" bson:"author"`
	Category         string              `json:"category" bson:"category"`
	Comments         []*comments.Comment `json:"comments" bson:"comments"`
	Created          string              `json:"created" bson:"created"`
	Id               string              `json:"id" bson:"id,omitempty"`
	Score            int                 `json:"score" bson:"score"`
	Text             string              `json:"text" bson:"text"`
	URL              string              `json:"url" bson:"url"`
	Title            string              `json:"title" bson:"title"`
	Type             string              `json:"type" bson:"type"`
	UpvotePercentage int                 `json:"upvotePercentage" bson:"upvotePercentage"`
	Views            int                 `json:"views" bson:"views"`
	Votes            []*VoteInfo         `json:"votes" bson:"votes"`
}

type VoteInfo struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

type NewPostForm struct {
	Category string            `json:"category"`
	Text     string            `json:"text"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`
	Url      string            `json:"url"`
	Author   *session.UserInfo `json: author`
}

type PostsHandler struct {
	PostsRepo PostsRepository
	UsersRepo userHandler.UserRepository
	Session   userHandler.SessionRepository
}

type PostsRepository interface {
	Get(id string) (*Post, error)
	GetAll() ([]Post, error)
	Add(*NewPostForm) (*Post, error)
	GetAllRelCategories(category string) ([]Post, error)
	GetAllRelUser(author *session.UserInfo) ([]Post, error)
	Upvote(id string, author *session.UserInfo) (*Post, error)
	DownVote(id string, author *session.UserInfo) (*Post, error)
	Unvote(id string, author *session.UserInfo) (*Post, error)
	Delete(id string) error
	AddComment(id string, comment *comments.Comment) (*Post, error)
	DeleteComment(id string, messageId string) (*Post, error)
}

func (p *PostsHandler) ListAll(w http.ResponseWriter, r *http.Request) {

	posts, err := p.PostsRepo.GetAll()
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, posts)
}

func (p *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	fd := &NewPostForm{}
	if err := utils.GetDataFromBody(w, r, fd); err != nil {
		return
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	fd.Author = sess
	newPost, err := p.PostsRepo.Add(fd)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, newPost)
}

func (p *PostsHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["POST_ID"]
	postData, err := p.PostsRepo.Get(postId)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, postData)
}

func (p *PostsHandler) ListRefCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["CATEGORY_NAME"]
	if len(category) == 0 {
		jsonError.NewJsonError(w, http.StatusInternalServerError, "there are no categories")
		return
	}
	posts, err := p.PostsRepo.GetAllRelCategories(category)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, posts)
}

func (p *PostsHandler) ListRefUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["USER_LOGIN"]
	if len(user) == 0 {
		return
	}
	author, err := p.UsersRepo.Get(user)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	session := &session.UserInfo{Username: author.Username, Id: author.ID}
	posts, err := p.PostsRepo.GetAllRelUser(session)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, posts)
}

func (p *PostsHandler) Remove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["POST_ID"]
	if len(id) != 0 {
		err := p.PostsRepo.Delete(id)
		if err != nil {
			jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (p *PostsHandler) VoteUp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["POST_ID"]

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := p.PostsRepo.Upvote(id, sess)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, post)
}

func (p *PostsHandler) VoteDown(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["POST_ID"]

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	post, err := p.PostsRepo.DownVote(id, sess)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, post)
}

func (p *PostsHandler) RemoveVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["POST_ID"]

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := p.PostsRepo.Unvote(id, sess)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, post)
}
