package main

import (
	"fmt"
	"strings"
)

func findAnagrams(words []string) map[string][]string {
	anagramGroups := make(map[string][]string)
	firstEntrances := make(map[[33]int]string)

	for _, w := range words {
		word := strings.ToLower(w)

		key := [33]int{}
		for _, r := range word {
			key[r-'а'] += 1
		}

		if firstEntrance, ok := firstEntrances[key]; ok {
			if _, ok := anagramGroups[firstEntrance]; ok {
				anagramGroups[firstEntrance] = append(anagramGroups[firstEntrance],  word)
			} else {
				anagramGroups[firstEntrance] = []string{firstEntrance, word}
			}
		} else {
			firstEntrances[key] = word
		}
	}

	return anagramGroups
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	fmt.Println(findAnagrams(input))
}
