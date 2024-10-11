package extractor

import (
	"bytes"
	"dev09/utils"
	"fmt"
	"golang.org/x/net/html"
	"path/filepath"
	"strings"
)

// CorrectAndExtractHTMLLinks Корректируем ссылки из html страницы и сохраняем их в мапу.
func CorrectAndExtractHTMLLinks(url string, node *html.Node, links *map[string]bool) {
	prevDir, err := utils.GetPathToRoot(url)
	if err == nil {
		for i, a := range node.Attr {
			if a.Key == "href" || a.Key == "src" || a.Key == "style" { // Проверяем по ключу, что этот атрибут содержит ссылку.
				if !strings.HasPrefix(a.Val, "#") { // Игнорирование якорных ссылок.
					if a.Key == "style" && strings.Contains(a.Val, "background-image:") {
						node.Attr[i].Val = strings.Replace(node.Attr[i].Val, "'", "", 2)
						(*links)[extractLinkFromAttrURL(node.Attr[i].Val)] = true
						node.Attr[i].Val = strings.Replace(node.Attr[i].Val, "/", prevDir, 1)
					} else {
						(*links)[a.Val] = true // Добавляем по значению ссылки в нашу мапу.
						node.Attr[i].Val = utils.NormalizeURL(a.Val)
					}
					if filepath.Ext(a.Val) != ".html" && a.Key != "style" {
						node.Attr[i].Val = prevDir + node.Attr[i].Val
					}
				}
			}
		}
	} else {
		fmt.Printf("failed to get path to root %v: %v", url, err)
	}

	for c := (*node).FirstChild; c != nil; c = c.NextSibling { // Рекурсивный цикл, который ходит по дочерним узлам дерева.
		CorrectAndExtractHTMLLinks(url, c, links) // Передаем сюда следующую ноду.
	}
}

// ExtractLinksFromCSS Достаем ссылки на файлы из css файла.
func ExtractLinksFromCSS(links *map[string]bool, buf *bytes.Buffer) {
	cssStr := buf.String() // Приводим наш css файл из буфера в строку.

	for val := range cssStr {
		if cssStr[val] == 'u' && cssStr[val+1] == 'r' && cssStr[val+2] == 'l' {
			ur := extractLinkFromAttrURL(cssStr[val:]) // Получаем ссылку из css файла.
			if !strings.Contains(ur, "data:") {        // Для игнорирования Data URI
				(*links)[ur] = true
			}
		}
	}
}

// Получаем ссылки в url().
func extractLinkFromAttrURL(dirPath string) string {
	start := strings.Index(dirPath, "url(") // Получаем индекс начала ссылки.
	end := strings.Index(dirPath, ")")      // Получаем индекс конца ссылки.
	sharpIdx := strings.Index(dirPath, "#") // В css указывает на определенные элементы графики внутри файла.

	if start == -1 || end == -1 { // Если не нашли индекс, то выходим.
		return ""
	}

	if sharpIdx != -1 && sharpIdx < end { // Если у нас нашелся шарп и он раньше скобки, то берем путь до шарпа.
		return dirPath[start+4 : sharpIdx]
	} // Иначе берем путь до скобки.

	return dirPath[start+4 : end]
}
