package getText

import (
	"log/slog"
	"net/http"
)

type getText interface {
	getText(alias string) (string, error)
}

func TextGetHandler(log *slog.Logger, getText getText) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.text.getText"

		log = log.With(
			slog.String("op", op),
		)

		//TODO: написать микросервис который бы отдавал текст по alias
	}
}
