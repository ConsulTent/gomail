package gomail

import (
	"log"
	"testing"
)

func TestHTMLMail(t *testing.T) {
	sender := &Sender{
		User:   "",
		Passwd: "",
		Host:   "guerrillamail.com",
		Port:   25,
	}
	sender.Configure()
	sw := sender.NewSendWorker(
		"fromuser@domain.local",
		"pumpkins@sharklasers.com",
		"testing",
	)

	err := sw.ParseTemplate(
		`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
        "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>

</head>

<body>
<p>
    Hello {{.Name}},
    Please visit us at <a href="{{.URL}}">our website.</a>
</p>

</body>

</html>`,
		map[string]string{
			"Name": "John Smith",
			"URL": "http://consultent.ltd/",
		},
	)

	if err != nil {
		log.Panicln(err)
	}

	err = sw.SendEmail()

	if err != nil {
		log.Panicln(err)
	}

	log.Println("mail module test passing")
}
