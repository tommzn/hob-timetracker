package timetracker

import (
	"os"
	"strings"
)

// NewFilePublisher returns a new publisher which writes report content to files.
func NewFilePublisher(path *string) *FilePublisher {
	filePath := "./"
	if path != nil {
		filePath = *path
	}
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	return &FilePublisher{FileMode: 0644, Path: filePath}
}

// FilePublisher wirtes contents to files.
type FilePublisher struct {
	FileMode os.FileMode
	Path     string
}

// Send will write passed content to given file name.
func (publisher *FilePublisher) Send(content []byte, fileName string) error {
	return os.WriteFile(publisher.Path+fileName, content, publisher.FileMode)
}
