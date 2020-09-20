package hw03_frequency_analysis //nolint:golint,stylecheck
import (
	"regexp"
	"sort"
	"strings"
)

func Top10(input string) []string {
	separatorRe := regexp.MustCompile(`[.,?!;:\s']+`)
	nonWordRe := regexp.MustCompile(`^[-—]+$|^$`)

	words := separatorRe.Split(input, -1)
	stats := make(map[string]int)
	for k, word := range words {
		if nonWordRe.MatchString(word) {
			continue
		}

		word = strings.ToLower(word)
		words[k] = word

		stats[word]++
	}

	maxSize := 10
	curSize := 0
	result := make([]string, 0, maxSize)
	min := -1
	for word, stat := range stats {
		changed := false
		if curSize < maxSize {
			result = append(result, word)
			changed = true
			curSize++
		} else if stat > min {
			result[maxSize-1] = word
			changed = true
		}

		if changed {
			sort.Slice(result, func(i, j int) bool {
				word1 := result[i]
				word2 := result[j]
				return stats[word1] > stats[word2]
			})
			minWord := result[curSize-1]
			min = stats[minWord]
		}
	}

	return result
}
