package timetracker

import "os"

// NewFilePublisher returns a new publisher which writes report content to files.
func NewFilePublisher() *FilePublisher {
	return &FilePublisher{FileMode: 0644}
}

// FilePublisher wirtes contents to files.
type FilePublisher struct {
	FileMode os.FileMode
}

// Send will write passed content to given file name.
func (publisher *FilePublisher) Send(content []byte, fileName string) error {
	return os.WriteFile(fileName, content, publisher.FileMode)
}
