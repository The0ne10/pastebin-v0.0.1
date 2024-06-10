package getText

import (
	resp "app/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Text string `json:"text"`
}

type getText interface {
	FindUUID(alias string) (string, error)
}

type getFile interface {
	ReadFile(alias string) (string, error)
}

func ReadTextHandler(log *slog.Logger, getText getText, file getFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.text.getText"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		uid, err := getText.FindUUID(alias)
		if err != nil {
			render.JSON(w, r, resp.Error("Not found"))
			log.Error("Not found alias", slog.String("error", err.Error()))

			return
		}

		fileRead, err := file.ReadFile(uid)
		if err != nil {
			render.JSON(w, r, resp.Error("Not found"))
			log.Error("Not found file", slog.String("error", err.Error()))

			return
		}

		render.JSON(w, r, fileRead)

		log.Info("File successfully read", slog.String("uuid", uid))
	}
}
