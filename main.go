package main

import (
	"bufio"
	"fmt"
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
	lines := make([]string, 0, 1024)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	st, err := logen.NewSanitizer()
	if err != nil {
		panic(err)
	}

	fmt.Println("sanitizing")
	sanitizedLines := make([]string, 0, len(lines))
	for i := range lines {
		sanitizedLines = append(sanitizedLines, st.Sanitized(lines[i]))
	}
	fmt.Println("matching")

	classes := make([]Match, 0, 1024)

	for i := 0; i < len(lines); i++ {
		var matched bool
		for j := range classes {
			matchLine := classes[j].event
			if matchEqual(sanitizedLines[i], sanitizedLines[matchLine]) {
				matched = true

				classes[j].occurrences = append(classes[j].occurrences, i)

				slices.SortFunc(classes, matchRank)

				break
			}
		}

		if !matched {
			classes = append(classes, Match{
				event:       i,
				occurrences: []int{i},
			})
		}

	}

	for _, class := range classes {
		fmt.Printf(`%d matches with event "%s" for line "%s"`, len(class.occurrences), strcase.ToCamel(makeEvent(sanitizedLines[class.event])), lines[class.event])
		fmt.Println()
	}

}

type Match struct {
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
