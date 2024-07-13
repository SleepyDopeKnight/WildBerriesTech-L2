package main

/*
Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
на входе подается файл с несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно:

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов
*/

import (
	f "dev03/internal/file"
	fls "dev03/internal/flags"
	s "dev03/internal/sort"
	"fmt"
	"log"
	"os"
)

func main() {
	fileLines, err := f.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalln(err)
	}

	fl := fls.FlagParse()
	s.NewSort(&fileLines, fl)
	
	for _, str := range fileLines {
		fmt.Println(str)
	}
}
