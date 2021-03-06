package notify

import (
	"bytes"
	"context"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

/**
	Modified from https://gist.github.com/jpillora/cb46d183eca0710d909a
	Thank you very much.
**/

const (
	/**
		Gmail SMTP Server
	**/
	SMTPServer = "smtp.gmail.com"
)

// Sender object it`s used to send SMTP request
type Sender struct {
	// User ...
	User string
	// Password
	Password string
}

// NewSender constructor
func NewSender(Username, Password string) Sender {

	return Sender{Username, Password}
}

// SendMail send mail to an array of dst
func (sender Sender) SendMail(ctx context.Context, Dest []string, Subject, bodyMessage string) {
	_, span := otel.Tracer(nameNotifierSMTP).Start(ctx, "SendMail")
	defer span.End()
	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, Dest, []byte(msg))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	span.SetAttributes(attribute.String("notifier.sendmail.client", "success"))

}

// WriteEmail an email
func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = sender.User

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	header["To"] = receipient
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

// WriteHTMLEmail ...
func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/html", subject, bodyMessage)
}

// WritePlainEmail ...
func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/plain", subject, bodyMessage)
}
