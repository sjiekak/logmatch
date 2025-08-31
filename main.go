package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"

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

	for _, class := range res.classes {
		fmt.Printf(`%d matches with event "%s" for line "%s"`, len(class.occurrences), class.key, res.lines[class.event])
		fmt.Println()
	}
}

func logmatch(r io.Reader) (state, error) {
	st, err := logen.NewSanitizer()
	if err != nil {
		return state{}, err
	}

	lines := make([]string, 0, 1024)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return state{}, err
	}

	fmt.Println("sanitizing")
	sanitizedLines := make([]string, 0, len(lines))
	for i := range lines {
		sanitizedLines = append(sanitizedLines, st.Sanitized(lines[i]))
	}
	fmt.Println("matching")

	// let's keep this ordered
	// the class which matched the most lines is assumed to be the most likely to match the next lines
	classes := make([]Match, 0, 1024)

	for i := 0; i < len(lines); i++ {
		var matched bool
		for j := range classes {
			matchLine := classes[j].event
			if matchEqual(sanitizedLines[i], sanitizedLines[matchLine]) {
				matched = true

				classes[j].occurrences = append(classes[j].occurrences, i)

				// sort when out of order
				// could also find desired location and swap but lazy
				if j != 0 && len(classes[j-1].occurrences) < len(classes[j].occurrences) {
					slices.SortFunc(classes, matchRank)
				}

				break
			}
		}

		if !matched {
			classes = append(classes, Match{
				key:         strcase.ToCamel(makeEvent(sanitizedLines[i])),
				event:       i,
				occurrences: []int{i},
			})
		}
	}

	return state{
		lines:          lines,
		sanitizedLines: sanitizedLines,
		classes:        classes,
	}, nil
}

type state struct {
	lines          []string
	sanitizedLines []string
	classes        []Match
}

type Match struct {
	key         string
	event       int
	occurrences []int
}

func matchRank(a, b Match) int {
	diff := len(b.occurrences) - len(a.occurrences)
	if diff != 0 {
		return diff
	}

	return b.event - a.event
}

func makeEvent(s string) string {
	loc := rePunctuation.FindStringIndex(s)
	if len(loc) == 0 {
		return s
	}

	return s[:loc[0]]
}
