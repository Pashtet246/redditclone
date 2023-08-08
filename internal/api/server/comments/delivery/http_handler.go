package commentsHandler

import (
	"net/http"
	"redditclone/internal/api/server/comments"
	postsRepo "redditclone/internal/api/server/posts/repo"
	"redditclone/internal/api/server/session"
	userHandler "redditclone/internal/api/server/users/delivery"
	"redditclone/utils"
	jsonError "redditclone/utils"

	"github.com/gorilla/mux"
)

type CommentsHandler struct {
	PostsRepo    *postsRepo.PostsRepository
	Session      userHandler.SessionRepository
	CommentsRepo comments.CommentRepository
}

func (c *CommentsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["POST_ID"]

	fd := &comments.CommentForm{}
	if err := utils.GetDataFromBody(w, r, fd); err != nil {
		return
	}
	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	comment, errAdding := c.CommentsRepo.Add(fd.Comment, sess)
	if errAdding != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, errAdding.Error())
		return
	}
	post, err := c.PostsRepo.AddComment(postId, comment)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, post)
}

func (c *CommentsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postId := vars["POST_ID"]
	commentId := vars["COMMENT_ID"]

	sess, err := session.SessionFromContext(r.Context())
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	errRemove := c.CommentsRepo.Delete(commentId, sess)
	if errRemove != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, errRemove.Error())
		return
	}
	post, err := c.PostsRepo.DeleteComment(postId, commentId)
	if err != nil {
		jsonError.NewJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.PrepareDataForSend(w, r, post)
}
