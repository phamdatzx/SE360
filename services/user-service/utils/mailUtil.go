package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail gửi email, cấu hình lấy từ biến môi trường:
// EMAIL_HOST, EMAIL_PORT, EMAIL_USERNAME, EMAIL_PASSWORD
func SendEmail(to []string, subject, body string) error {
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")

	if host == "" || port == "" || username == "" || password == "" {
		return fmt.Errorf("email configuration missing in environment variables")
	}

	// Tạo message
	message := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		username, to[0], subject, body))

	// Kết nối tới SMTP server
	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Gửi email
	if err := smtp.SendMail(addr, auth, username, to, message); err != nil {
		return err
	}

	return nil
}

// BuildActivationEmailContent tạo nội dung email kích hoạt tài khoản
func BuildActivationEmailContent(token, email string) string {
	clientUrl := os.Getenv("CLIENT_URL")
	activationLink := fmt.Sprintf(clientUrl+"/activate?token=%s", token)

	content := fmt.Sprintf(`
Hi %s,

Cảm ơn bạn đã đăng ký tài khoản.

Vui lòng nhấn vào link dưới đây để kích hoạt tài khoản của bạn:

%s

Link này sẽ hết hạn sau 24 giờ.

Trân trọng,
Your App Team
`, email, activationLink)

	return content
}

func BuildResetPasswordEmailContent(token, email string) string {
	clientUrl := os.Getenv("CLIENT_URL")
	activationLink := fmt.Sprintf(clientUrl+"/reset-password?token=%s", token)

	content := fmt.Sprintf(`
Hi %s,

Cảm ơn bạn đã đăng ký tài khoản.

Vui lòng nhấn vào link dưới đây để kích hoạt tài khoản của bạn:

%s

Link này sẽ hết hạn sau 24 giờ.

Trân trọng,
Your App Team
`, email, activationLink)

	return content
}
