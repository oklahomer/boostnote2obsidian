package main

import (
	"boostnote2obsidian/boostnote"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	from := flag.String("from", "", "Path to Boost Note file directory")
	to := flag.String("to", "", "Directory path to export converted files")
	flag.Parse()

	if *from == "" || *to == "" {
		log.Fatal("./path/main.go -from=/path/to/boostnote/storage -to=/path/to/target")
	}

	err := os.MkdirAll(*to, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to setup target directory")
	}

	err = filepath.Walk(*from, func(filePath string, info os.FileInfo, e error) error {
		if !strings.HasSuffix(filePath, ".json") {
			// Skip an irrelevant file
			return nil
		}

		note, e := boostnote.ReadNote(filePath)
		if e != nil {
			log.Printf("Failed to read note %s: %s\n", filePath, e.Error())
			return nil
		}

		targetPath := path.Join(*to, genFileName(note))
		file, e := os.Create(targetPath)
		if e != nil {
			return fmt.Errorf("failed to create file: %w", e)
		}
		defer func() {
			file.Close()
		}()

		_, e = file.WriteString(note.Content)
		if e != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, e)
		}
		return nil
	})
	if err != nil {
		log.Fatal("%s", err.Error())
	}
}

var safeTitle = regexp.MustCompile(`[^[:alnum:]-.]`)
var sepChars = regexp.MustCompile(`[ &_=+:]`)
var dashChars = regexp.MustCompile(`[\-]+`)

func genFileName(note *boostnote.Note) string {
	name := strings.ToLower(note.Title)
	name = strings.TrimSpace(name)
	name = sepChars.ReplaceAllString(name, "-")
	name = safeTitle.ReplaceAllString(name, "")  // Non ascii to empty string
	name = dashChars.ReplaceAllString(name, "-") // Multiple dashes to single

	if name == "" {
		name = fmt.Sprintf(
			"%d-%02d-%02d-%02d-%02d",
			note.CreatedAt.Year(),
			note.CreatedAt.Month(),
			note.CreatedAt.Day(),
			note.CreatedAt.Hour(),
			note.CreatedAt.Minute())
	}

	return fmt.Sprintf("%s.md", name)
}
