package createFile

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

func Create(name uuid.UUID, data []byte) (string, error) {
	const op = "lib.createFile.Create"

	currentDir, _ := os.Getwd()
	path := filepath.Join(currentDir, "files/")

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	path = filepath.Join(path, name.String()+".txt")

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}

	if _, err = file.Write(data); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	return path, nil
}
