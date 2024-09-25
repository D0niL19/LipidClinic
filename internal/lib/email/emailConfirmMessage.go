package email

func EmailConfirmMessage(url string) []byte {
	return []byte("Subject: " + "Подтверждение регистрации" + "\r\n" +
		"\r\n" +
		"Пожалуйста, подтвердите свою регистрацию, перейдя по следующей ссылке." + "\r\n" +
		url)

}
