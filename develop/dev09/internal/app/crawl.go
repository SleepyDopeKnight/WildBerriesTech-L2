package app

import (
	"dev09/internal/downloader"
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
			if link != "/" { // Записывает иногда ссылки на другие сайты в таком виде, чтобы их пропускать тоже.
				absPAth, err := utils.ConvertRelativeURLToAbsolute(url, link)
				if err != nil {
					fmt.Println("Failed to convert link to absolute: ", err)
				} else {
					Crawl(absPAth, rootDomain, visited)
				}
			}
		}
	}
}
