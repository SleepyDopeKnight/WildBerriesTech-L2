package utils

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// GetFullPath Получает полный путь, от корневой директории сайта.
func GetFullPath(ur string) (string, error) {
	u, err := url.Parse(ur)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL to get full path: %s", err)
	}

	return u.Host + u.Path, nil
}

// GetPathToRoot Получает путь, который вернет к корневой директории сайта.
func GetPathToRoot(ur string) (string, error) {
	u, err := url.Parse(ur)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL to get root path: %s", err)
	}

	fullPath, err := GetFullPath(ur)
	if err != nil {
		return "", fmt.Errorf("GetFullPath error in GetPathToRoot: %s", err)
	}

	fp := strings.Trim(fullPath, u.Host)     // Получаем путь и вырезаем рутовую директорию.
	fpLen := len(strings.Split(fp, "/")) - 2 // -2 т.к. в начале пути есть слэш.

	return strings.Repeat("../", fpLen), nil
}

// GetFileName Получает имя для файла вместе с путем до него и расширением, для рутовой дает имя index.html.
func GetFileName(dirPath string) string {
	if strings.Trim(dirPath, path.Dir(dirPath)) == "/" { // Если остается слэш, то получаем имя главной страницы index.html
		return dirPath + "index.html"
	}

	return SetExtHTML(dirPath) // Проверяем нужно ли добавить расширение файлу.
}

// SetExtHTML Устанавливает расширение .html для html страниц, если его нет.
func SetExtHTML(url string) string {
	if filepath.Ext(path.Base(url)) == "" { // Проверяем содержит ли имя файла расширение.
		url += ".html" // Если не содержит, то добавляем .html
	}

	return url
}
