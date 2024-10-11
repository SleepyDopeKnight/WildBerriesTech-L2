package creator

import (
	"dev09/utils"
	"fmt"
	"os"
	"path"
	"strings"
)

// CreateFile Создает файл и нужную для него директорию в зависимости от пути.
func CreateFile(ur string) (*os.File, error) {
	err := createDir(ur) // Создаем директории.
	if err != nil {
		return nil, err
	}
	fullPath, err := utils.GetFullPath(ur) // Получаем путь созданных директорий.
	if err != nil {
		return nil, fmt.Errorf("GetFullPath error in CreateFile: %v", err)
	}

	file, err := os.Create(utils.GetFileName(fullPath)) // Создаем сам файл в созданных директориях.
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %s", err)
	}

	return file, nil
}

func createDir(ur string) error {
	fullPath, err := utils.GetFullPath(ur) // Получаем путь созданных директорий.

	if err != nil {
		return fmt.Errorf("GetFullPath error in createDir: %v", err)
	}

	// Отрезаем у пути имя файла и создаем директории.
	err = os.MkdirAll(strings.TrimSuffix(fullPath, "/"+path.Base(fullPath)), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed creating directory: %s", err)
	}
	return nil
}
