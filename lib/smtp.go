package lib

import (
	"context"
	"fmt"
	"net/smtp"
	"strconv"

	"gitlab.com/tsmdev/software-development/backend/go-project/dto"
)

type SMTP struct {
	Context  context.Context
	Logger   Logger
	Sender   string
	Server   string
	Port     int
	Username string
	Password string
}

func NewSMTP(env Env, logger Logger) SMTP {
	return SMTP{
		Context:  context.Background(),
		Logger:   logger,
		Username: env.SMTPUser,
		Password: env.SMTPPass,
		Server:   env.SMTPHost,
		Port:     env.SMTPPort,
		Sender:   env.SMTPSenderName,
	}
}

func (s SMTP) SendEmail(request dto.SendEmailRequestDto) error {

	from := fmt.Sprintf("From: <%s>\r\n", s.Sender)
	to := fmt.Sprintf("To: <%s>\r\n", request.To)
	subject := "Subject: " + request.Subject + "\r\n"
	body := request.Body + "\r\n"

	message := from + to + subject + "\r\n" + body

	s.Logger.Info("Email message: ", message)

	err := smtp.SendMail(s.Server+":"+strconv.Itoa(s.Port),
		smtp.PlainAuth("", s.Username, s.Password, s.Server),
		s.Sender, []string{request.To}, []byte(message),
	)

	if err != nil {
		s.Logger.Error("Error when sending email")
		s.Logger.Debug("Detail: ", err.Error())
		return err
	}

	return nil
}
