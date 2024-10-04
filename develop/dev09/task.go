package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
var (
	firstIn    bool   = true                  // Нужна для создания index.html, после первого входа она всегда false.
	visited           = make(map[string]bool) // Сюда записываем посещенные урлы, чтобы рекурсия не пошла по кругу.
	rootDomain string                         // Нужна для проверки, что рекурсия не ушла дальше заданного ресура.
)

func main() {
	u, err := url.Parse("https://www.marinohall.com")
	if err != nil {
		fmt.Println(err)
	}

	rootDomain = u.Host
	rootLink := u.Scheme + "://" + u.Host + "/" // Для того чтобы обход начинался с корня, а подать ссылку мог не корня
	crawl(rootLink, rootDomain)
}

func crawl(url, rootDomain string) {
	if visited[url] { // Если посещали уже, то выходим из функции.
		return
	}

	visited[url] = true                                                   // Записываем урл в мапу посещенных.
	if strings.Contains(url, rootDomain) && filepath.Ext(url) != ".dmg" { // Смотрим, содрежит ли урл изначальный домен, чтобы рекурсия не пошла дальше.
		fmt.Println("Downloading: ", url)

		links, err := downloadAndExtractLinks(url) // Скачиваем страницу и достаем ссылки.
		if err != nil {
			fmt.Println("Failed to download page: ", err)
		}

		for link := range links {
			if link != "/" { // Записывает иногда ссылки на другие в таком виде, чтобы их пропускать тоже.
				crawl(convertRelativeUrlToAbsolute(url, link), rootDomain)
			}
		}
	}
}

func extractLinksFromNodeAndCorrectPath(url string, node *html.Node, links *map[string]bool) {
	fp := strings.Trim(createDirAndGetPath(url), rootDomain) // Получаем путь и вырезаем рутовую директорию.
	fpSlice := strings.Split(fp, "/")
	fpLen := len(fpSlice) - 2 // -2 т.к. в начале пути есть слэш.
	prevDir := strings.Repeat("../", fpLen)

	for i, a := range node.Attr {
		if a.Key == "href" || a.Key == "src" || a.Key == "style" { // Проверяем по ключу, что этот атрибут содержит ссылку.
			if !strings.HasPrefix(a.Val, "#") { // Игнорирование якорных ссылок.
				(*links)[a.Val] = true // Добавляем по значению ссылки в нашу мапу.
				node.Attr[i].Val = TrimSlashAndInitHome(a.Val)
				if a.Key == "style" && strings.Contains(a.Val, "background-image:") {
					node.Attr[i].Val = strings.Replace(node.Attr[i].Val, "/", prevDir, 1)
					node.Attr[i].Val = strings.Replace(node.Attr[i].Val, "'", "", 2)
				}
				if filepath.Ext(a.Val) != ".html" && a.Key != "style" {
					node.Attr[i].Val = prevDir + node.Attr[i].Val
				}
			}
		}
	}

	for c := (*node).FirstChild; c != nil; c = c.NextSibling { // Рекурсивный цикл, который ходит по дочерним узлам дерева.
		extractLinksFromNodeAndCorrectPath(url, c, links) // Передаем сюда следующую ноду.
	}
}

func downloadAndExtractLinks(url string) (map[string]bool, error) {
	buf, err := downloadPage(url) // Скачиваем страницу и сохраняем тело в буфер.
	if err != nil {
		return nil, err
	}

	file, err := createFile(url) // Создаем файл.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	links, err := extractLinksAndSaveFile(url, buf, file)
	if err != nil {
		return nil, err
	}

	return links, nil
}

func downloadPage(url string) (*bytes.Buffer, error) {
	resp, err := http.Get(url) // C помощью get запроса получаем содержимое страницы.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request: %s", resp.Status)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body) // Сохраняем данные в буфер, для переиспользования body.
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func extractLinksAndSaveFile(url string, buf *bytes.Buffer, file *os.File) (map[string]bool, error) {
	links := make(map[string]bool) // Сюда будем записывать ссылки со страницы.

	doc, _ := html.Parse(bytes.NewReader(buf.Bytes()))   // Парсим нашу страницу из буфера через ридер, для повторного чтения.
	extractLinksFromNodeAndCorrectPath(url, doc, &links) // Извлекаем ссылки на другие страницы.

	if path.Ext(file.Name()) == ".html" {
		err := html.Render(file, doc) // Используем рендер для записи, тк io.copy не сохранит изменение ссылок внутри страницы.
		if err != nil {
			return nil, fmt.Errorf("error render to file: %s", err)
		}
	} else {
		_, err := io.Copy(file, bytes.NewReader(buf.Bytes())) // Копируем данные в файл из буфера через ридер, для повторного чтения.
		if err != nil {
			return nil, fmt.Errorf("error writing to file: %s", err)
		}
	}
	return links, nil
}

func createFile(url string) (*os.File, error) {
	dp := createDirAndGetPath(url)

	file, err := os.Create(getFileName(dp)) // Создаем сам файл.
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %s", err)
	}

	return file, nil
}

func createDirAndGetPath(ur string) string {
	u, err := url.Parse(ur)
	if err != nil {
		fmt.Println("Failed to parse URL: ", err)
		return ""
	}

	err = os.MkdirAll(u.Host+path.Dir(u.Path), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	return u.Host + u.Path
}

func convertRelativeUrlToAbsolute(pageUrl, href string) string {
	base, err := url.Parse(pageUrl)
	if err != nil {
		fmt.Println("failed to parse base url: ", err)
		return ""
	}

	ref, err := url.Parse(href) // href может быть как /about, так и example.com/about.
	if err != nil {
		fmt.Println("failed to parse href: ", err)
		return ""
	}

	return base.ResolveReference(ref).String() // Получаем путь, склеиваем или оставляем неизменным в зависимости от href.
}

//func convertToLocalPath(href string) string {
//
//}

func TrimSlashAndInitHome(href string) string {
	if href == "/" {
		return "index.html"
	} else {
		return setExtHTML(strings.TrimPrefix(href, "/"))
	}
}

func setExtHTML(url string) string {
	if filepath.Ext(path.Base(url)) == "" { // Проверяем содержит ли имя файла расширение.
		url += ".html" // Если не содержит, то добавляем .html
	}

	return url
}

func getFileName(dp string) string {
	if firstIn { // Если первый вызов, то получаем имя главной страницы index.html
		firstIn = false
		return path.Dir(dp) + "/index.html"
	}
	return setExtHTML(dp) // Проверяем нужно ли добавить расширение файлу.
}
