package timetracker

import "os"

// NewFilePublisher returns a new publisher which writes report content to files.
func NewFilePublisher(path *string) *FilePublisher {
	if path == nil {
		localDir := "./"
		path = &localDir
	}
	return &FilePublisher{FileMode: 0644, Path: *path}
}

// FilePublisher wirtes contents to files.
type FilePublisher struct {
	FileMode os.FileMode
	Path     string
}

// Send will write passed content to given file name.
func (publisher *FilePublisher) Send(content []byte, fileName string) error {
	return os.WriteFile(fileName, content, publisher.FileMode)
}
