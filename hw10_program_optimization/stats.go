package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, tld string) (DomainStat, error) {
	result := DomainStat{}
	if tld == "" {
		return result, nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	re := regexp.MustCompile("@(.+\\." + tld + ")")
	for scanner.Scan() {
		line := scanner.Bytes()

		email := jsoniter.Get(line, "Email").ToString()

		matches := re.FindStringSubmatch(strings.ToLower(email))
		if len(matches) == 0 {
			continue
		}

		domain := matches[1]
		result[domain]++
	}

	return result, nil
}
