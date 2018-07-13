package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/standupdev/strset"
)

func main() {
	if len(os.Args) > 1 {
		query := strings.Join(os.Args[1:], " ")
		data, err := os.Open("UnicodeData.txt")
		if err != nil {
			panic(err)
		}
		defer data.Close()
		linesResult := make(chan string)
		go filter(data, query, linesResult)
		for line := range linesResult {
			fmt.Println(line)
		}
	} else {
		fmt.Println("Please enter one or more words to search.")

	}
}

func parseLine(line string) (string, rune) {
	fields := strings.Split(line, ";")
	code, _ := strconv.ParseInt(fields[0], 16, 32)
	return fields[1], rune(code)
}

func match(query strset.Set, name string) bool {
	nameTerms := strset.MakeFromText(name)
	return query.SubsetOf(nameTerms)
}

func matcher(query strset.Set) func(name string) bool {
	queryMatch := query
	return func(name string) bool {
		return match(queryMatch, name)
	}
}

func filter(data io.Reader, query string, lines chan<- string) {
	queryTerms := strset.MakeFromText(strings.ToUpper(query))
	scanner := bufio.NewScanner(data)
	myMatcher := matcher(queryTerms)
	for scanner.Scan() {
		name, code := parseLine(scanner.Text())
		if myMatcher(name) {
			lines <- fmt.Sprintf("U+%04X\t%c\t%s", code, code, name)
		}
	}
	close(lines)
}
