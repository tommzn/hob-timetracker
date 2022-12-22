package timetracker

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FilePublisherTestSuite struct {
	suite.Suite
}

func TestFilePublisherTestSuite(t *testing.T) {
	suite.Run(t, new(FilePublisherTestSuite))
}

func (suite *FilePublisherTestSuite) TestSend() {

	publisher := NewFilePublisher(nil)
	fileName := "test.file"
	suite.Nil(publisher.Send([]byte("test"), fileName))
	_, err := os.Stat(fileName)
	suite.False(errors.Is(err, os.ErrNotExist))

	os.Remove(fileName)
}
