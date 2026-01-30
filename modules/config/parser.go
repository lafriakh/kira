package config

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"sync"
)

// Parse the file and return a map of configs.
func Parse(src io.Reader) (map[string]any, error) {
	var data = make(map[string]any)
	var lock = sync.RWMutex{}

	lock.Lock()
	defer lock.Unlock()

	// Start reading from the reader using a scanner.
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip the line if we start with hashtag.
		if len(line) > 0 && string(line[0]) == "#" {
			continue
		}

		before := beforeDelimiter(line, delimiter)
		after := afterDelimiter(line, delimiter)

		if before == "" || after == "" {
			continue
		}

		afterString := after

		if len(afterString) >= 2 && afterString[:1] == `"` && afterString[len(afterString)-1:] == `"` { // String
			data[before[:]] = afterString[1 : len(afterString)-1]
		} else if toInt, err := strconv.ParseInt(afterString, 10, 64); err == nil { // Number
			data[before[:]] = toInt
		} else if toFloat, err := strconv.ParseFloat(afterString, 64); err == nil { // float
			data[before[:]] = toFloat
		} else if toBool, err := strconv.ParseBool(afterString); err == nil { // bool
			data[before[:]] = toBool
		} else {
			data[before[:]] = after
		}
	}

	if scanner.Err() != nil {
		return data, scanner.Err()
	}

	return data, nil
}

func beforeDelimiter(value string, a string) string {
	before, _, ok := strings.Cut(value, a)
	if !ok {
		return ""
	}
	return strings.TrimSpace(before)
}

func afterDelimiter(value string, a string) string {
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return strings.TrimSpace(value[pos+1:])
}
