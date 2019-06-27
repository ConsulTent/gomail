package gomail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/smtp"
)

// Sender global Sender config
type Sender struct {
	User   string
	Passwd string
	Host   string
	Port   int
	auth   smtp.Auth
}

// SendWorker worker process
type SendWorker struct {
	*Sender
	from    string
	to      string
	subject string
	body    string
	head    string
}

// NewSendWorker return a initialized SendWorker struct
func (s *Sender) NewSendWorker(from, to, subject string) *SendWorker {
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s", from)
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"
	var message string
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	return &SendWorker{
		Sender:  s,
		from:    from,
		to:      to,
		subject: subject,
		head:    message,
	}
}

// Configure make smtp plain auth
func (s *Sender) Configure() {
	if len(s.Passwd) > 0 {
	s.auth = smtp.PlainAuth(
		"",
		s.User,
		s.Passwd,
		s.Host,
	)
  }
	return
}

func (sw *SendWorker) ParseTemplate(tmpl string, data interface{}) error {
	t, err := template.New("eMailTemplate").Parse(tmpl)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	sw.body = buffer.String()
	return nil
}

// SendEmail send email
func (sw *SendWorker) SendEmail() error {
	url := fmt.Sprintf("%s:%d", sw.Sender.Host, sw.Port)

	err := smtp.SendMail(
		url,
		sw.auth,
		sw.from,
		[]string{sw.to},
		[]byte(sw.head+sw.body),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// Dial return a smtp client
func Dial(addr string) (*smtp.Client, error) {

	config := tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", addr, &config)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
