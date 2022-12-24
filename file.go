package timetracker

import (
	"os"
	"strings"

	log "github.com/tommzn/go-log"
)

// NewFilePublisher returns a new publisher which writes report content to files.
func NewFilePublisher(path *string, logger log.Logger) *FilePublisher {
	filePath := "./"
	if path != nil {
		filePath = *path
	}
	if !strings.HasSuffix(filePath, "/") {
		filePath += "/"
	}
	return &FilePublisher{FileMode: 0644, Path: filePath, logger: logger}
}

// FilePublisher wirtes contents to files.
type FilePublisher struct {
	FileMode os.FileMode
	Path     string
	logger   log.Logger
}

// Send will write passed content to given file name.
func (publisher *FilePublisher) Send(content []byte, fileName string) error {
	if err := os.WriteFile(publisher.Path+fileName, content, publisher.FileMode); err != nil {
		publisher.logger.Error("Unable to write report to file, reason: ", err)
		return err
	}
	publisher.logger.Debug("Report successful written to file: ", publisher.Path+fileName)
	return nil
}
