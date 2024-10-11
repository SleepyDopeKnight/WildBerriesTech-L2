package downloader

import (
	"bytes"
	cr "dev09/internal/creator"
	ex "dev09/internal/extractor"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// DownloadAndExtractLinks Скачивает страницу и сохраняет с нее все ссылки в мапу.
func DownloadAndExtractLinks(url string) (map[string]bool, error) {
	buf, err := downloadPage(url) // Скачиваем страницу и сохраняем тело в буфер.
	if err != nil {
		return nil, err
	}

	file, err := cr.CreateFile(url) // Создаем файл.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	links, err := extractLinksAndSaveFile(url, buf, file) // Извлекаем ссылки и сохраняем страницу в файл.
	if err != nil {
		return nil, err
	}

	if filepath.Ext(url) == ".css" {
		ex.ExtractLinksFromCSS(&links, buf) // Если css файл, то ищем в нем ссылки дизайна сайта.
	}

	return links, nil
}

func downloadPage(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url) // C помощью get запроса получаем содержимое страницы.
	if err != nil {
		return nil, fmt.Errorf("http get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request: %s", resp.Status)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body) // Сохраняем данные в буфер, для переиспользования body.
	if err != nil {
		return nil, fmt.Errorf("failed copy body to buffer: %s", err)
	}

	return &buf, nil
}

func extractLinksAndSaveFile(url string, buf *bytes.Buffer, file *os.File) (map[string]bool, error) {
	links := make(map[string]bool) // Сюда будем записывать ссылки со страницы.

	doc, err := html.Parse(bytes.NewReader(buf.Bytes())) // Парсим нашу страницу из буфера через ридер, для повторного чтения.
	if err != nil {
		return nil, fmt.Errorf("html parse: %v", err)
	}

	ex.CorrectAndExtractHTMLLinks(url, doc, &links) // Извлекаем ссылки на другие страницы.

	if path.Ext(file.Name()) == ".html" {
		err = html.Render(file, doc) // Используем рендер для записи, тк io.copy не сохранит изменение ссылок внутри страницы.
		if err != nil {
			return nil, fmt.Errorf("error render to file: %s", err)
		}
	} else {
		_, err = io.Copy(file, bytes.NewReader(buf.Bytes())) // Копируем данные в файл из буфера через ридер, для повторного чтения.
		if err != nil {
			return nil, fmt.Errorf("error writing to file: %s", err)
		}
	}

	return links, nil
}
