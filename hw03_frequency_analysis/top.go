package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`([а-я]+[\-а-я]*)+`)

type Word struct {
	Value     string
	Frequency int
}

func NewWord(value string) *Word {
	return &Word{
		Value:     value,
		Frequency: 1,
	}
}

type Words []Word

func (ws Words) Sort() {
	sort.Slice(ws, func(i, j int) bool {
		if ws[i].Frequency == ws[j].Frequency {
			return ws[i].Value < ws[j].Value
		}
		return ws[i].Frequency > ws[j].Frequency
	})
}

func (ws Words) ToStrings(count int) []string {
	result := make([]string, 0, count)
	for i := range ws {
		result = append(result, ws[i].Value)
		if i == count-1 {
			break
		}
	}
	return result
}

type WordsMap map[string]*Word

func (wm WordsMap) ToSlice() Words {
	result := make(Words, 0, len(wm))
	for i := range wm {
		result = append(result, *wm[i])
	}
	return result
}

func calcWordsFrequency(in []string) WordsMap {
	m := make(WordsMap, len(in))
	for _, str := range in {
		if m[str] == nil {
			m[str] = NewWord(str)
		} else {
			m[str].Frequency++
		}
	}
	return m
}

func Top10(in string) []string {
	words := reg.FindAllString(strings.ToLower(in), -1)
	wordsWithRate := calcWordsFrequency(words)
	wordsSlice := wordsWithRate.ToSlice()
	wordsSlice.Sort()
	return wordsSlice.ToStrings(10)
}
