package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ed-henrique/somos-dev-br/internal/models"
	"github.com/ed-henrique/somos-dev-br/internal/ui/pages"
)

func GetAllPosts(q *models.Queries) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts, err := q.GetPosts(r.Context())
		if err != nil {
			errMessage := "could not get all posts"
			slog.Error(errMessage, slog.String("handler", "GetAllPosts"), slog.String("err", err.Error()))
			http.Error(w, errMessage, http.StatusInternalServerError)
			return
		}

		err = pages.Posts(posts).Render(r.Context(), w)
		if err != nil {
			errMessage := "could not render 'get all posts'"
			slog.Error(errMessage, slog.String("handler", "GetAllPosts"), slog.String("err", err.Error()))
			http.Error(w, errMessage, http.StatusInternalServerError)
		}
	})
}
