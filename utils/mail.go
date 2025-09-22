package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sonit_server/constant/env"
	"sonit_server/constant/noti"
	"sonit_server/model/dto/request"
	"text/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendMail(req request.SendMailRequest) error {
	var errLogMsg string = fmt.Sprintf(noti.MAIL_ERR_MSG, "Util.Mail - SendMail")

	template, err := template.ParseFiles(req.TemplatePath)
	if err != nil {
		req.Logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, req.Body); err != nil {
		req.Logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	// serviceEmail := os.Getenv(env.SERVICE_EMAIL)
	// securityPass := os.Getenv(env.SECURITY_PASS)
	// host := os.Getenv(env.HOST)
	// port, err := strconv.Atoi(os.Getenv(env.MAIL_PORT))
	// if err != nil {
	// 	port = 587
	// }

	// m := mail.NewMessage()
	// m.SetHeader("From", "taikhoanhoconl123@gmail.com")
	// m.SetHeader("To", req.Body.Email)
	// m.SetHeader("Subject", req.Body.Subject)
	// m.SetBody("text/html", body.String())

	// dialer := mail.NewDialer(host, port, serviceEmail, securityPass)

	// if err := dialer.DialAndSend(m); err != nil {
	// 	req.Logger.Println(errLogMsg + err.Error())
	// 	return errors.New(noti.GENERATE_MAIL_WARN_MSG)
	// }

	res, err := sendgrid.NewSendClient(os.Getenv(env.SONIT_MAIL_KEY)).Send(mail.NewSingleEmail(
		mail.NewEmail("Sonit Custom", "phuchtqse183980@fpt.edu.vn"),
		req.Body.Subject,
		mail.NewEmail("", req.Body.Email),
		"",
		body.String(),
	))

	if err != nil {
		req.Logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	if res.StatusCode >= 400 {
		req.Logger.Println(errLogMsg + res.Body)
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}
