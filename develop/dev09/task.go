package main

import (
	"fmt"
	"golang.org/x/net/html"
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
var visited = make(map[string]bool)
var firstIn bool = true

func main() {

	u, err := url.Parse("https://www.marinohall.com")
	if err != nil {
		fmt.Println(err)
	}

	rootLink := u.Scheme + "://" + u.Host + "/" // Для того чтобы обход начинался с корня, а подать ссылку мог не корня
	crawl(rootLink)
}

func crawl(url string) {
	if visited[url] {
		return
	}
	visited[url] = true

	fmt.Println("Downloading: ", url)
	if err := downloadPage(url); err != nil {
		fmt.Println("Failed to download page: ", err)
	}

	links := extractLinks(url)
	for _, link := range links {
		crawl(convertRelativeUrlToAbsolute(url, link))
	}
}

func extractLinks(url string) []string {
	var links []string // Сюда будем записывать ссылки со страницы.

	resp, err := http.Get(url) // C помощью get запроса получаем содержимое страницы.
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body) // Парсим нашу страницу.
	if err != nil {
		fmt.Println(err)
		return nil
	}

	extractLinksFromNode(doc, &links)

	return links
}

func extractLinksFromNode(node *html.Node, links *[]string) {
	if node.Type == html.ElementNode && node.Data == "a" { // Проверяем что элемент ноды является тегом со ссылкой.
		for _, a := range node.Attr {
			if a.Key == "href" { // Проверяем по ключу, что этот атрибут содержит ссылку.
				if !strings.HasPrefix(a.Val, "#") { // Игнорирование якорных ссылок.
					*links = append(*links, a.Val) // Добавляем по значению ссылки в наш слайс.
				}
			}
		}
	}

	for c := (*node).FirstChild; c != nil; c = c.NextSibling { // Рекурсивный цикл, который ходит по дочерним узлам дерева.
		extractLinksFromNode(c, links) // Передаем сюда следующую ноду.
	}
}

func jopa(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" { // Проверяем что элемент ноды является тегом со ссылкой.
		for i, a := range node.Attr {
			if a.Key == "href" { // Проверяем по ключу, что этот атрибут содержит ссылку.
				if !strings.HasPrefix(a.Val, "#") { // Игнорирование якорных ссылок.
					node.Attr[i].Val = convertToLocalPath(a.Val)
				}
			}
		}
	}

	for c := (*node).FirstChild; c != nil; c = c.NextSibling { // Рекурсивный цикл, который ходит по дочерним узлам дерева.
		jopa(c) // Передаем сюда следующую ноду.
	}
}

func downloadPage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad request: %s", resp.Status)
	}

	var fileName string
	if firstIn {
		fileName = "index.html"
	} else {
		fileName = path.Base(url)
		if !strings.Contains(fileName, ".") {
			fileName += ".html"
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//_, err = io.Copy(file, resp.Body)
	doc, _ := html.Parse(resp.Body)
	jopa(doc)
	html.Render(file, doc)
	//return fmt.Errorf("error writing HTML to file: %s", err)
	firstIn = false

	return err
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
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	var fileName string
	if strings.Contains(path.Base(href), ".") {
		fileName = path.Base(href)
	} else {
		fileName = path.Base(href + ".html")
	}

	return wd + "/" + fileName
}
