package email

import (
	"LipidClinic/internal/config"
	"net/smtp"
)

func SendEmail(recipient string, message []byte, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.Source, cfg.Smtp.Password, cfg.Smtp.Host)

	err := smtp.SendMail(cfg.Smtp.Host+":"+cfg.Smtp.Port, auth, cfg.Source, []string{recipient}, message)
	if err != nil {
		return err
	}

	return nil
}
