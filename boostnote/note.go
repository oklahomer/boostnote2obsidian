package boostnote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var ErrTrashedNote = errors.New("given note is already trashed")

func ReadNote(path string) (*Note, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read note: %w", err)
	}

	file := &boostNoteFile{}
	err = json.Unmarshal(bytes, file)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize note content: %w", err)
	}
	if file.Trashed {
		return nil, ErrTrashedNote
	}

	firstLineBreakAt := strings.Index(file.Content, "\n")
	if firstLineBreakAt < 0 {
		return nil, errors.New("no title found")
	}
	title := []byte(file.Content)[:firstLineBreakAt]
	content := []byte(file.Content)[firstLineBreakAt+1:]

	note := &Note{
		CreatedAt:      file.CreatedAt,
		UpdatedAt:      file.UpdatedAt,
		FolderPathName: file.FolderPathName,
		Title:          string(title),
		Content:        strings.TrimSpace(string(content)),
		Tags:           file.Tags,
	}
	return note, nil
}

type Note struct {
	CreatedAt      *NoteTime
	UpdatedAt      *NoteTime
	FolderPathName string
	Title          string
	Content        string
	Tags           []string
}

type boostNoteFile struct {
	CreatedAt      *NoteTime `json:"createdAt"`
	UpdatedAt      *NoteTime `json:"updatedAt"`
	FolderPathName string    `json:"folderPathname"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Tags           []string  `json:"tags"`
	Trashed        bool      `json:"trashed"`
}

type NoteTime struct {
	time.Time
}

const (
	NoteTimeFormat = "2006-01-02T15:04:05.999Z"
)

func (t *NoteTime) UnmarshalText(data []byte) error {
	tm, err := time.Parse(`"`+NoteTimeFormat+`"`, string(data))
	t.Time = tm
	return err
}
