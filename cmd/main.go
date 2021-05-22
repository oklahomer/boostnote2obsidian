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
	"strconv"
	"strings"
)

func main() {
	from := flag.String("from", "", "Path to Boost Note file directory")
	to := flag.String("to", "", "Directory path to export converted files")
	flag.Parse()

	// Check required params
	if *from == "" || *to == "" {
		log.Fatal("./path/main.go -from=/path/to/boostnote/storage -to=/path/to/target")
	}

	err := filepath.Walk(*from, func(filePath string, info os.FileInfo, e error) error {
		if !strings.HasSuffix(filePath, ".json") {
			// Skip an irrelevant file
			return nil
		}

		// Read note file
		note, e := boostnote.ReadNote(filePath)
		if e != nil {
			log.Printf("Failed to read note %s: %s\n", filePath, e.Error())
			return nil
		}

		// Prepare export destination
		targetPath := genFilePath(*to, note)
		dir := path.Dir(targetPath)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", e)
		}
		file, e := os.Create(targetPath)
		if e != nil {
			return fmt.Errorf("failed to create file: %w", e)
		}
		defer func() {
			file.Close()
		}()

		// Export content
		_, e = file.WriteString(note.Content)
		if e != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, e)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk through the directory: %s", err.Error())
	}
}

var safeTitle = regexp.MustCompile(`[^[:alnum:]-.]`)
var sepChars = regexp.MustCompile(`[ &_=+:]`)
var dashChars = regexp.MustCompile(`[\-]+`)

func genFilePath(baseDir string, note *boostnote.Note) string {
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
	name = fmt.Sprintf("%s.md", name)

	return path.Join(baseDir, strconv.Itoa(note.CreatedAt.Year()), fmt.Sprintf("%02d", note.CreatedAt.Month()), name)
}
