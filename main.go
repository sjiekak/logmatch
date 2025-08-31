package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/iancoleman/strcase"
	"github.com/sjiekak/logen"
)

var (
	rePunctuation = regexp.MustCompile("[[:punct:]]")
)

func matchEqual(a string, b string) bool {
	return a == b
}

func main() {
	res, err := logmatch(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, group := range res {
		fmt.Printf(`%d matches with event "%s" for line "%s"`, len(group.occurrences), group.eventName, group.representative)
		fmt.Println()
	}
}

func logmatch(r io.Reader) ([]Match, error) {
	st, err := logen.NewSanitizer()
	if err != nil {
		return nil, err
	}

	// let's keep this ordered
	// the class which matched the most lines is assumed to be the most likely to match the next lines
	matchGroups := make([]Match, 0, 16)
	scanner := bufio.NewScanner(r)

	lineNum := 0
	for scanner.Scan() {
		currentLine := scanner.Text()
		currentSanitized := st.Sanitized(currentLine)

		var matched bool
		for i := range matchGroups {
			if matchEqual(matchGroups[i].representative, currentSanitized) {
				matched = true

				matchGroups[i].occurrences = append(matchGroups[i].occurrences, lineNum)
				moveToCorrectRank(matchGroups, i)

				break
			}
		}

		if !matched {
			matchGroups = append(matchGroups, Match{
				eventName:      makeEvent(currentSanitized),
				representative: currentSanitized,
				rawLine:        currentLine,
				occurrences:    []int{lineNum},
			})
		}

		lineNum++
	}

	return matchGroups, scanner.Err()
}

// move the element at index in its correct position in the sorted match group list
// works only because the index is already >= index + n
func moveToCorrectRank(matchGroups []Match, index int) {
	j := index - 1
	for j >= 0 && len(matchGroups[j].occurrences) < len(matchGroups[index].occurrences) {
		j--
	}

	j++
	if j == index {
		return
	}
	matchGroups[j], matchGroups[index] = matchGroups[index], matchGroups[j]
}

type state struct {
	lines          []string
	sanitizedLines []string
	classes        []Match
}

type Match struct {
	eventName      string
	representative string
	rawLine        string
	occurrences    []int
}

func makeEvent(s string) string {
	event := s
	loc := rePunctuation.FindStringIndex(s)
	if len(loc) > 0 {
		event = s[:loc[0]]
	}

	return strcase.ToCamel(event)
}
