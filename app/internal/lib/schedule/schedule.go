package schedule

import (
	"log/slog"
	"time"
)

type findFilesDeletedAt interface {
	DeleteExpired() ([]string, error)
}

type DeletedExpiredFiles interface {
	DeleteFiles(log *slog.Logger, uids []string) error
}

func StartSchedule(
	log *slog.Logger,
	uid findFilesDeletedAt,
	files DeletedExpiredFiles,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	log.Info("starting schedule")
	go func() {
		for {
			select {
			case <-ticker.C:
				uids, err := uid.DeleteExpired()
				if err != nil {
					log.Error("Error deleting expired records:", slog.String("error", err.Error()))
				}
				if err = files.DeleteFiles(log, uids); err != nil {
					log.Error("Error deleting files:", slog.String("error", err.Error()))
				}
			}
		}
	}()
}
