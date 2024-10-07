package main

import (
	"dev09/internal/app"
	"fmt"
	"net/url"
)

/*
=== Утилита wget ===

# Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var rootDomain string // Нужна для проверки, что рекурсия не ушла дальше заданного ресура.

func main() {
	u, err := url.Parse("https://www.marinohall.com")
	if err != nil {
		fmt.Println(err)
	}

	visited := make(map[string]bool) // Сюда записываем посещенные урлы, чтобы рекурсия не пошла по кругу.
	rootDomain = u.Host
	rootLink := u.Scheme + "://" + u.Host + "/" // Для того чтобы обход начинался с корня, а подать ссылку мог не корня
	app.Crawl(rootLink, rootDomain, &visited)
}
