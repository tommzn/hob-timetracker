package timetracker

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

// NewEMailPublisher creates a new publsher to send time tracking reports via email.
func NewEMailPublisher(source, destination, subject, message string) *EMailPublisher {
	return &EMailPublisher{
		Source:      source,
		Destination: destination,
		Subject:     subject,
		Message:     message,
	}
}

// EMailPublisher delivers time tracking reports via email.
type EMailPublisher struct {
	Source, Destination, Subject, Message string
}

// Send will deliver given time tracking report via email.
func (publisher *EMailPublisher) Send(content []byte, fileName string) error {

	rawEMail, err := rawEMail(publisher.Source, publisher.Destination, publisher.Subject, publisher.Message, content, fileName)
	if err != nil {
		return err
	}

	rawMessage := ses.RawMessage{
		Data: []byte(*rawEMail),
	}
	sendRawEmailInput := &ses.SendRawEmailInput{
		Destinations: []*string{aws.String(publisher.Destination)},
		Source:       aws.String(publisher.Source),
		RawMessage:   &rawMessage,
	}

	client := newSESClient(nil)
	_, sendErr := client.SendRawEmail(sendRawEmailInput)
	return sendErr
}

// RawEMail generates a raw email with given sender/receiver, subject, message and attachment.
func rawEMail(source, destination, subject, message string, csvFile []byte, attachmentFilename string) (*string, error) {

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	h := make(textproto.MIMEHeader)
	h.Set("From", source)
	h.Set("To", destination)
	h.Set("Return-Path", source)
	h.Set("Subject", subject)
	h.Set("Content-Language", "en-US")
	h.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")
	h.Set("MIME-Version", "1.0")
	_, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}

	h = make(textproto.MIMEHeader)
	h.Set("Content-Transfer-Encoding", "7bit")
	h.Set("Content-Type", "text/html; charset=us-ascii")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write([]byte(message))
	if err != nil {
		return nil, err
	}

	fn := attachmentFilename
	h = make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "attachment; filename="+fn)
	h.Set("Content-Type", "text/csv; x-unix-mode=0644; name=\""+fn+"\"")
	h.Set("Content-Transfer-Encoding", "7bit")
	part, err = writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(csvFile)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	s := buf.String()
	if strings.Count(s, "\n") < 2 {
		return nil, fmt.Errorf("invalid e-mail content")
	}
	s = strings.SplitN(s, "\n", 2)[1]
	return &s, nil
}
