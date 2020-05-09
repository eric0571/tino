package config

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/smtp"

	"io/ioutil"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

// EmailInfo is the details for the SMTP server
type EmailInfo struct {
	Username string
	Password string
	Hostname string
	Port     int
	From     string
}

// SendEmail sends an email
func (e EmailInfo) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", e.Username, e.Password, e.Hostname)
	header := make(map[string]string)
	header["From"] = e.From
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = `text/html; charset="utf-8"`
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	// Send the email
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", e.Hostname, e.Port),
		auth,
		e.From,
		[]string{to},
		[]byte(message),
	)
	if err != nil {
		fmt.Println("Sen Mail error:", err)
	}
	return err
}

type Literal interface {
	io.Reader
	Len() int
}

type ResponseMessage struct {
	From    string
	Subject string
	Content string
	Date    string
}

func AcceptAllMail(addr, user, pass string) ([]Literal, error) {
	client, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}
	defer client.Logout()
	if err := client.Login(user, pass); err != nil {
		log.Fatal(err)
	}
	// 收件箱
	mbox, err := client.Select("INBOX", true)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		fmt.Println("收件箱中没有邮件")
		return nil, nil
	}
	fmt.Println("InBox message count:", mbox.Messages)

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), mbox.Messages)
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []imap.FetchItem{"BODY[]"}, messages)
		fmt.Println("Message Count:", len(messages))
	}()
	// 返回体
	response := []Literal{}
	// 收件箱的所有邮件
	for msg := range messages {
		section, err := imap.ParseBodySectionName("BODY[]")
		if err != nil {
			return nil, fmt.Errorf("邮件解析错误")
		}

		// r := msg.GetBody("BODY[]")
		r := msg.GetBody(section)
		if r == nil {
			return nil, fmt.Errorf("没有邮件内容")
		}
		response = append(response, r)
	}
	return response, nil
}

func GetMailList() []*ResponseMessage {
	// 返回体
	response := []*ResponseMessage{}
	request, err := AcceptAllMail("imap.163.com:993", "wlh3354@163.com", "Xykzdykd5")
	if err != nil {
		log.Print(err)
		return response
	}

	for _, r := range request {
		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			fmt.Println("Create mail reader err:", err)
			return nil
		}

		var existEntity = new(ResponseMessage)
		header := mr.Header
		if from, err := header.AddressList("From"); err == nil {
			for _, value := range from {
				existEntity.From = value.Address
			}
		}

		if subject, err := header.Subject(); err == nil {
			existEntity.Subject = subject
		}

		if date, err := header.Date(); err == nil {
			existEntity.Date = date.Format("2006-01-02 15:04:05")
		}

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Println("Err:", err)
				// log.Fatal(err)
			}
			switch p.Header.(type) {
			// case mail.TextHeader:
			default:
				// This is the message's text (can be plain-text or HTML)
				b, err := ioutil.ReadAll(p.Body)
				if err != nil {
					fmt.Println(err)
				}
				existEntity.Content = string(b)
				fmt.Println("Mail content:", existEntity.Content)
			}
		}
		response = append(response, existEntity)
	}

	return response
}
