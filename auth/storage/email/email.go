package email

import (
	"crypto/tls"
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func SendEmail(email string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("email.myEmail"))
	m.SetHeader("To", email)
	m.SetHeader("New_user", "Registration notice")
	m.SetBody("text/html", "<p><a href='dog.html'> Собаки </a></p>")
	d := gomail.NewDialer("smtp.gmail.com", 587, "anvladislav59@gmail.com", viper.GetString("email.password"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return errors.Errorf("Не удалось отправить сообщение о регистрации:%s", err)
	}

	return nil
}
