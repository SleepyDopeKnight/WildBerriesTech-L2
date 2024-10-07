package utils

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// GetFullPath Получает полный путь, от корневой директории сайта.
func GetFullPath(ur string) string {
	u, err := url.Parse(ur)
	if err != nil {
		fmt.Println("Failed to parse URL to get path: ", err)
		return ""
	}

	return u.Host + u.Path
}

// GetPathToRoot Получает путь, который вернет к корневой директории сайта.
func GetPathToRoot(ur string) string {
	u, err := url.Parse(ur)
	if err != nil {
		fmt.Println("Failed to parse URL to get root path: ", err)
		return ""
	}

	fp := strings.Trim(GetFullPath(ur), u.Host) // Получаем путь и вырезаем рутовую директорию.
	fpLen := len(strings.Split(fp, "/")) - 2    // -2 т.к. в начале пути есть слэш.

	return strings.Repeat("../", fpLen)
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
