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
	createDir(ur)                     // Создаем директории.
	fullPath := utils.GetFullPath(ur) // Получаем путь созданных директорий.

	file, err := os.Create(utils.GetFileName(fullPath)) // Создаем сам файл в созданных директориях.
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %s", err)
	}

	return file, nil
}

func createDir(ur string) {
	fullPath := utils.GetFullPath(ur) // Получаем путь созданных директорий.

	// Отрезаем у пути имя файла и создаем директории.
	err := os.MkdirAll(strings.TrimSuffix(fullPath, "/"+path.Base(fullPath)), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}
