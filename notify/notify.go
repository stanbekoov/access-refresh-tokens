package notify

import (
	"fmt"
	"medods-test/db"
	"net/smtp"
	"os"
)

func SendIpAddressWarn(uid string, ip string) error {
	user := db.User{Uid: uid}

	r := db.Db.First(&user)
	if r.RowsAffected == 0 {
		return fmt.Errorf("SendIpAddressWarn: no user with this id found")
	}

	if ip == user.Ip {
		return nil
	}

	sender := os.Getenv("EMAIL_SENDER")
	pass := os.Getenv("EMAIL_SENDER_PASS")
	host := os.Getenv("EMAIL_SENDER_HOST")
	port := os.Getenv("EMAIL_PORT")

	receiver := user.Email

	auth := smtp.PlainAuth("", sender, pass, host)

	msg := []byte("Someone has loginned to your account from unknown IP address")

	err := smtp.SendMail(host+":"+port, auth, sender, []string{receiver}, msg)
	if err != nil {
		return err
	}

	return nil
}
