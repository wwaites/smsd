package main

import (
	"crypto/tls"
	"net/smtp"
)

func InsecureSendMail(domain string, addr string, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(domain); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: domain, InsecureSkipVerify: true}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

func SendMail(cfg Config, rt Route, m Message) (err error) {
	to   := rt.Dst
	from := m.Src + "@" + cfg.Server.Domain
	bmsg := []byte("To: " + to + "\r\n" +
		"From:  " + from + "\r\n" +
		"Subject: SMS Message\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		m.Msg)

	err = InsecureSendMail(cfg.Server.Domain, cfg.Server.Mta, from, []string{to}, bmsg)
	return
}
