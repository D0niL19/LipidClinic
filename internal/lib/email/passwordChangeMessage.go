package email

func PasswordChangeMessage(url string) []byte {
	return []byte("Subject: " + "Восстановление пароля" + "\r\n" +
		"\r\n" +
		"Перейдите по следующей ссылке для восстановления пароля:" + "\r\n" +
		url)

}
