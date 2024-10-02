package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/
var firstIn bool = true             // Нужна для создания index.html, после первого входа она всегда false.
var visited = make(map[string]bool) // Сюда записываем посещенные урлы, чтобы рекурсия не пошла по кругу.

func main() {

	u, err := url.Parse("https://www.marinohall.com")
	if err != nil {
		fmt.Println(err)
	}

	rootDomain := u.Host                        // Нужна для проверки, что рекурсия не ушла дальше заданного ресура.
	rootLink := u.Scheme + "://" + u.Host + "/" // Для того чтобы обход начинался с корня, а подать ссылку мог не корня
	crawl(rootLink, rootDomain)
}

func crawl(url, rootDomain string) {
	if visited[url] { // Если посещали уже, то выходим из функции.
		return
	}

	visited[url] = true                    // Записываем урл в мапу посещенных.
	if strings.Contains(url, rootDomain) { // Смотрим, содрежит ли урл изначальный домен, чтобы рекурсия не пошла дальше.
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

func extractLinksFromNode(node *html.Node, links *map[string]bool) {
	if node.Type == html.ElementNode && node.Data == "a" { // Проверяем что элемент ноды является тегом со ссылкой.
		for i, a := range node.Attr {
			if a.Key == "href" { // Проверяем по ключу, что этот атрибут содержит ссылку.
				if !strings.HasPrefix(a.Val, "#") { // Игнорирование якорных ссылок.
					(*links)[a.Val] = true // Добавляем по значению ссылки в нашу мапу.
					node.Attr[i].Val = convertToLocalPath(a.Val)
				}
			}
		}
	}

	for c := (*node).FirstChild; c != nil; c = c.NextSibling { // Рекурсивный цикл, который ходит по дочерним узлам дерева.
		extractLinksFromNode(c, links) // Передаем сюда следующую ноду.
	}
}

func downloadAndExtractLinks(url string) (map[string]bool, error) {
	links := make(map[string]bool) // Сюда будем записывать ссылки со страницы.

	resp, err := http.Get(url) // C помощью get запроса получаем содержимое страницы.
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad request: %s", resp.Status)
	}

	file, err := createFile(url) // Создаем файл.
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc, _ := html.Parse(resp.Body)   // Парсим нашу страницу.
	extractLinksFromNode(doc, &links) // Извлекаем ссылки на другие страницы.
	_, err = io.Copy(file, resp.Body) // Копируем страницу в файл.
	if err != nil {
		return nil, fmt.Errorf("error writing HTML to file: %s", err)
	}

	return links, nil
}

func createFile(url string) (*os.File, error) {
	var fileName string
	if firstIn { // Если первый вызов, то получаем имя главной страницы index.html
		fileName = "index.html"
	} else {
		fileName = getFileName(url) // Получаем имя файла.
	}

	file, err := os.Create(fileName) // Создаем сам файл.
	if err != nil {
		return nil, fmt.Errorf("failed creating file: %s", err)
	}

	firstIn = false // После первого вызова зануляем наш флаг.

	return file, nil
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

func convertToLocalPath(href string) string {
	var fileName string

	wd, err := os.Getwd() // Получаем путь, где мы находимся
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fileName = getFileName(href) // Получаем имя файла.

	return wd + "/" + fileName
}

func getFileName(url string) string {
	fileName := path.Base(url)
	if !strings.Contains(fileName, ".") { // Проверяем содержит ли имя файла расширение.
		fileName += ".html" // Если не содержит, то добавляем .html
	}

	return fileName
}
