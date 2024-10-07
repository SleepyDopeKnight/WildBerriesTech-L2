package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// ConvertRelativeURLToAbsolute конвертирует абсолютную ссылку в абсолютную.
func ConvertRelativeURLToAbsolute(pageURL, href string) string {
	base, err := url.Parse(pageURL)
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

// NormalizeURL если домашняя страница, то дает имя index.html, по надобности добавляет расширение .html и отрезает начальный слэш.
func NormalizeURL(href string) string {
	if href == "/" { // Если домашняя страница, то это главная страница.
		return "index.html"
	}
	return SetExtHTML(strings.TrimPrefix(href, "/")) // Если нужно, то добавляем html разрешение и отрезаем /
}
