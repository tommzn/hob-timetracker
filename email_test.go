package timetracker

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type EMailTestSuite struct {
	suite.Suite
}

func TestEMailTestSuite(t *testing.T) {
	suite.Run(t, new(EMailTestSuite))
}

func (suite *EMailTestSuite) skipCI() {
	if _, isSet := os.LookupEnv("CI"); isSet {
		suite.T().Skip("Skip test in CI environment.")
	}
}

func (suite *EMailTestSuite) TestSendEMail() {

	suite.skipCI()

	source := suite.emailSourceForTest()
	destination := suite.emailDestinationForTest()
	subject := "Time Tracking Report - TEST"
	message := "<html><h1>Time Tracking Report (TEST)</h1><p>Here's your rport!</p></html>"
	filename := "TestReport_202201.xlsx"

	fileContent, err := os.ReadFile("docs/" + filename)
	suite.Nil(err)

	publisher := NewEMailPublisher(source, destination, subject, message)
	suite.Nil(publisher.Send(fileContent, filename))
}

func (suite *EMailTestSuite) emailSourceForTest() string {
	source, ok := os.LookupEnv("HOB_EMAIL_SOURCE")
	suite.True(ok)
	return source
}

func (suite *EMailTestSuite) emailDestinationForTest() string {
	destination, ok := os.LookupEnv("HOB_EMAIL_DESTINATION")
	suite.True(ok)
	return destination
}
