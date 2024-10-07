package app

import (
	"dev09/downloader"
	"dev09/utils"
	"fmt"
	"path/filepath"
	"strings"
)

// Crawl Рекурсивно ходит по страницам, скачивает их и рекурсивно берет ссылки с этих страниц.
func Crawl(url, rootDomain string, visited *map[string]bool) {
	if (*visited)[url] { // Если посещали уже, то выходим из функции.
		return
	}

	(*visited)[url] = true                                                // Записываем урл в мапу посещенных.
	if strings.Contains(url, rootDomain) && filepath.Ext(url) != ".dmg" { // Смотрим, содрежит ли урл изначальный домен, чтобы рекурсия не пошла дальше.
		fmt.Println("Downloading: ", url)

		links, err := downloader.DownloadAndExtractLinks(url) // Скачиваем страницу и достаем ссылки.
		if err != nil {
			fmt.Println("Failed to download page: ", err)
		}

		for link := range links {
			if link != "/" { // Записывает иногда ссылки на другие в таком виде, чтобы их пропускать тоже.
				Crawl(utils.ConvertRelativeURLToAbsolute(url, link), rootDomain, visited)
			}
		}
	}
}
