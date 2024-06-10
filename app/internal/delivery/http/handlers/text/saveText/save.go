package saveText

import (
	resp "app/internal/lib/api/response"
	"app/internal/lib/createFile"
	"app/internal/lib/random"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"log/slog"
	"net/http"
	"time"
)

const (
	length = 7
)

type SaveText interface {
	Save(alias string, uid uuid.UUID, textTTL string) error
}

type SaveFile interface {
	UploadFile(uid uuid.UUID, path string) error
}

type Response struct {
	resp.Response
	Alias string `json:"url"`
}

type Request struct {
	Text    string `json:"text"`
	TimeTTL string `json:"time_ttl,omitempty"`
}

func SaveTextHandler(log *slog.Logger, saveText SaveText, file SaveFile) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.text.saveText"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to deserialize request", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode json"))
			return
		}

		// примитивная валидация
		// TODO: в следущей версии использовать (validator/v10)
		if req.Text == "" {
			render.JSON(w, r, resp.Error("text is required"))
			log.Error("text is required", slog.String("error", "text is required"))
			return
		}

		uid, err := uuid.NewRandom()
		if err != nil {
			log.Error("failed to generate uuid", slog.String("error", err.Error()))
			return
		}

		alias := random.GenerateRandomString(length)

		path, err := createFile.Create(uid, []byte(req.Text))
		if err != nil {
			log.Error("failed to create file", slog.String("error", err.Error()))
			return
		}
		err = file.UploadFile(uid, path)
		if err != nil {
			log.Error("failed to upload file", slog.String("error", err.Error()))
			render.JSON(w, r, "failed to upload file")
			return
		}

		if err != nil {
			log.Error("failed to create file", slog.String("error", err.Error()))
			render.JSON(w, r, "failed to create message")
			return
		}

		timeTTL, err := validationTimeTTL(req.TimeTTL)
		if err != nil {
			log.Error("failed to validate time ttl", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to validate time_ttl"))
			return
		}

		err = saveText.Save(alias, uid, timeTTL)
		if err != nil {
			log.Error("failed to save text", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to save text"))
			return
		}

		//TODO: идентифицировать пользователя если зарегистрирован

		responseSuccess(w, r, alias)

		log.Info("saved text", slog.String("alias", alias))
	}
}

func responseSuccess(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.Success(),
		Alias:    alias,
	})
}

func validationTimeTTL(interval string) (string, error) {
	// Если timeTTL пустое, установите время хранения по умолчанию на 24 часа
	if interval == "" {
		return time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"), nil
	}

	timeTTL, err := time.ParseDuration(interval)
	if err != nil {
		return "", fmt.Errorf("failed to parse time_ttl")
	}

	return time.Now().Add(timeTTL).Format("2006-01-02 15:04:05"), nil
}
