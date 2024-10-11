package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// ConvertRelativeURLToAbsolute конвертирует абсолютную ссылку в абсолютную.
func ConvertRelativeURLToAbsolute(pageURL, href string) (string, error) {
	base, err := url.Parse(pageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base url: %s", err)
	}

	ref, err := url.Parse(href) // href может быть как /about, так и example.com/about.
	if err != nil {
		return "", fmt.Errorf("failed to parse href: %s", err)
	}

	return base.ResolveReference(ref).String(), nil // Получаем путь, склеиваем или оставляем неизменным в зависимости от href.
}

// NormalizeURL если домашняя страница, то дает имя index.html, по надобности добавляет расширение .html и отрезает начальный слэш.
func NormalizeURL(href string) string {
	if href == "/" { // Если домашняя страница, то это главная страница.
		return "index.html"
	}
	return SetExtHTML(strings.TrimPrefix(href, "/")) // Если нужно, то добавляем html разрешение и отрезаем /
}
