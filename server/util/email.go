package util

import (
	"fmt"
	gomail "gopkg.in/mail.v2"
	"watt-flow/config"
)

func SendEmail(receiver string, subject string, body string) error {
	env := config.Init()
	message := gomail.NewMessage()
	message.SetHeader("From", "wattflow12@gmail.com")
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", env.EmailSecret)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	} else {
		return nil
	}
}

func GenerateActivationEmailBody(activationLink string) string {
	return fmt.Sprintf(`
    <html>
        <body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
            <div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
                <h2 style="color: #333; margin-bottom: 20px;">Welcome to Watt-Flow</h2>
                <p style="font-size: 16px; color: #555;">Thank you for registering!</p>
                <p style="font-size: 16px; color: #555;">Click the button below to activate your account:</p>
                <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Activate Account</a>
                <p style="font-size: 14px; color: #888; margin-top: 20px;">If you did not sign up for this account, please ignore this email.</p>
            </div>
        </body>
    </html>`, activationLink)
}

func GenerateSuccessfulActivationEmailBody(loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Account Successfully Activated</h2>
					<p style="font-size: 16px; color: #555;">You have successfully activated your account!</p>
					<p style="font-size: 16px; color: #555;">Click the button below to go to the login page:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Go to Login</a>
				</div>
			</body>
		</html>
	`, loginLink)
}
