package saveText

import (
	"log/slog"
	"net/http"
)

type SaveText interface {
	save(text string) (string, error)
}

func SaveTextHandler(log *slog.Logger, save SaveText) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http.handlers.text.saveText"

		log = log.With(
			slog.String("op", op),
		)

		//TODO: сделать уникальные сокращатели ссылок

		//TODO: сохранять url в базу данных

		//TODO: сделать возможность удаления пользователем ссылки на текст по времени (по умолчанию 5 минут)

		//TODO: идентифицировать пользователя
	}
}
