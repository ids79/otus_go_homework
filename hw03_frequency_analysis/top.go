package hw03frequencyanalysis

import (
	"math"
	"regexp"
	"sort"
	"strings"
)

type TopWords struct {
	word   string
	number int
}

var TaskWithAsteriskIsCompleted = false
var rC = regexp.MustCompile(`[,|.|!|?|"|:|;]`)

func getFunc(top []TopWords) func(i, j int) bool {
	return func(i, j int) bool {
		if top[i].number == top[j].number {
			return strings.Compare(top[i].word, top[j].word) == -1
		}
		return top[i].number > top[j].number
	}
}

func Top10(str string) []string {
	if TaskWithAsteriskIsCompleted {
		str = strings.ToLower(str)
		str = rC.ReplaceAllString(str, "")
	}
	words := strings.Fields(str)
	m := make(map[string]int)
	for _, word := range words {
		if TaskWithAsteriskIsCompleted && word == "-" {
			continue
		}
		m[word]++
	}
	top := make([]TopWords, 0)
	for word, n := range m {
		t := TopWords{word: word, number: n}
		top = append(top, t)
	}
	sort.Slice(top, getFunc(top))
	top10len := int(math.Min(float64(10), float64(len(top))))
	top10 := make([]string, top10len)
	for i := 0; i < 10; i++ {
		if len(top) <= i {
			break
		}
		top10[i] = top[i].word
	}

	return top10
}
