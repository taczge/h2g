package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
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

func toDate(line string) time.Time {
	date := strings.TrimPrefix(line, "BASENAME: ")
	format := "2006/01/02/150405"
	time, err := time.Parse(format, date)
	if err != nil {
		panic(err)
	}

	return time
}

type Entry struct {
	Id    string
	Title string
	Date  time.Time
	Body  string
}

func (e *Entry) WriteToFile(outdir string) {
	outpath := fmt.Sprintf("%s/%s", outdir, e.Id)
	f, err := os.Create(outpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString(fmt.Sprintf("TITLE: %s\n", e.Title))
	f.WriteString(fmt.Sprintf("DATE: %s\n", e.Date))
	f.WriteString(e.Body)
}

func split(filename, outdir string) {
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
			entry.WriteToFile(outdir)
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
	if len(os.Args) < 3 {
		fmt.Println("Useage: h2g INPUT_FILE OUT_DIR")
		os.Exit(1)
	}

	split(os.Args[1], os.Args[2])
}
