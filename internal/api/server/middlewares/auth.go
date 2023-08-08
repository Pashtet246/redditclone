package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"redditclone/internal/api/server/session"
	sessionRepo "redditclone/internal/api/server/session/repo"
)

var (
	NotRequiredAuthUrls = map[*regexp.Regexp]string{
		regexp.MustCompile(`/api/login`):    "POST",
		regexp.MustCompile(`/api/register`): "POST",
		regexp.MustCompile(`/api/posts$`):   "GET",
		regexp.MustCompile(`/api/post/\b[0-9a-f]{8}\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\b[0-9a-f]{12}\b`): "GET",
		regexp.MustCompile(`/api/posts`): "GET",
		regexp.MustCompile(`/api/user`):  "GET",
		regexp.MustCompile(`/$`):         "GET",
		regexp.MustCompile(`/static/`):   "GET",
		regexp.MustCompile(`/manifest`):  "GET",
		regexp.MustCompile(`/favicon`):   "GET",
	}
)

func Auth(sm *sessionRepo.SessionRepository, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for regexpPath, method := range NotRequiredAuthUrls {
			if regexpPath.MatchString(r.URL.String()) && method == r.Method {
				next.ServeHTTP(w, r)
				return
			}
		}
		inToken := r.Header.Get("Authorization")
		if inToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error while parsing header"))
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		splittedInToken := strings.Split(inToken, " ")
		if len(splittedInToken) < 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("error while parsing header - wrong type of token in header"))
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		token := splittedInToken[1]
		sess, err := sm.Check(token)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), session.SessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
