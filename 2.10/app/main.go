package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

func main() {
	// Парсинг флагов
	column := flag.Int("k", 0, "sort via column number")
	numeric := flag.Bool("n", false, "sort numerically")
	reverse := flag.Bool("r", false, "reverse sort order")
	unique := flag.Bool("u", false, "output only unique lines")
	sortByMonth := flag.Bool("M", false, "sort by month")
	ignoreTrailingBlanks := flag.Bool("b", false, "ignore trailing blanks")
	checkSorted := flag.Bool("c", false, "check if input is sorted")
	humanNumeric := flag.Bool("h", false, "sort by human-readable numbers")

	flag.Parse()

	app := &application{
		column:               *column,
		numeric:              *numeric,
		reverse:              *reverse,
		unique:               *unique,
		sortByMonth:          *sortByMonth,
		ignoreTrailingBlanks: *ignoreTrailingBlanks,
		checkSort:            *checkSorted,
		humanNumeric:         *humanNumeric,
	}

	// Чтение данных
	lines, err := readLines(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		os.Exit(1)
	}

	// Проверка сортировки
	if app.checkSort {
		if app.checkSorted(lines) {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// Подготовка ключей
	keys := make([]*key, len(lines))
	for i, line := range lines {
		keys[i] = app.key(line)
	}

	// Сортировка
	sort.Slice(keys, func(i, j int) bool {
		less := app.fullLess(keys[i], keys[j])
		if app.reverse {
			return !less
		}
		return less
	})

	// Обновление строк согласно отсортированным ключам
	for i := range lines {
		lines[i] = keys[i].original
	}

	// Удаление дубликатов
	if app.unique {
		_, lines = removeDuplicates(keys, lines)
	}

	// Вывод результатов
	for _, line := range lines {
		fmt.Println(line)
	}
}
