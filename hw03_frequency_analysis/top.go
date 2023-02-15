package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`([а-я]+[\-а-я]*)+`)

func calcWordsFrequency(in []string) map[string]int {
	m := make(map[string]int, len(in))
	for _, str := range in {
		m[str]++
	}
	return m
}

func Top10(in string) []string {
	words := reg.FindAllString(strings.ToLower(in), -1)
	wordsMap := calcWordsFrequency(words)

	wordsSlice := make([]string, 0, len(wordsMap))
	for i := range wordsMap {
		wordsSlice = append(wordsSlice, i)
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		valueI := wordsSlice[i]
		valueJ := wordsSlice[j]
		if wordsMap[valueI] == wordsMap[valueJ] {
			return wordsSlice[i] < wordsSlice[j]
		}
		return wordsMap[valueI] > wordsMap[valueJ]
	})

	limit := 10
	if len(wordsSlice) < 10 {
		limit = len(wordsSlice)
	}
	return wordsSlice[:limit]
}
