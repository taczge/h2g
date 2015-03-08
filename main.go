package main

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"time"
)

func toTitle(line string) string {
	title := strings.TrimPrefix(line, "TITLE: ")

	return strings.TrimSpace(title)
}

func toId(line string) string {
	name := strings.TrimPrefix(line, "BASENAME: ")

	return strings.Replace(name, "/", "-", -1)
}

var intToMonth = map[int]time.Month{
	1:  time.January,
	2:  time.February,
	3:  time.March,
	4:  time.April,
	5:  time.May,
	6:  time.June,
	7:  time.July,
	8:  time.August,
	9:  time.September,
	10: time.October,
	11: time.November,
	12: time.December,
}

func toDate(line string) time.Time {
	date := strings.TrimPrefix(line, "BASENAME: ")

	y, _ := strconv.Atoi(date[0:4])
	m, _ := strconv.Atoi(date[5:7])
	d, _ := strconv.Atoi(date[8:10])

	return time.Date(y, intToMonth[m], d, 0, 0, 0, 0, time.UTC)
}

type Entry struct {
	Id    string
	Title string
	Date  time.Time
	Body  string
}

func (e *Entry) WriteToFile() {
	f, err := os.Create(e.Id)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(e.Body)
}

func extractEntryData(filename string) {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	entry := Entry{}
	isInBody := false
	var body bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "BASENAME: ") {
			entry.Id = toId(line)
			entry.Date = toDate(line)
			continue
		}

		if strings.HasPrefix(line, "TITLE: ") {
			entry.Title = toTitle(line)
			continue
		}

		if strings.HasPrefix(line, "BODY:") {
			isInBody = true
			continue
		}

		if strings.HasPrefix(line, "--------") {
			entry.Body = body.String()
			entry.WriteToFile()
			body.Reset()
			isInBody = false
			continue
		}

		if strings.HasPrefix(line, "-----") {
			continue
		}

		if isInBody {
			body.WriteString(strings.TrimSpace(line))
			body.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	extractEntryData("n91.hatenablog.com.export.txt")
}