package util

import (
	"fmt"
	"watt-flow/config"

	gomail "gopkg.in/mail.v2"
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

func GenerateMonthlyBillEmail(
	billID uint64, billDate string, spentPower float64, price float64,
	ownerName string, validFrom string, blueZone, redZone, greenZone, billingPower, tax float64,
	paymentLink string,
) string {
	return fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif; background: #f4f4f4; color: #333; padding: 40px; text-align: center;">
			<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: left;">
				<h2 style="color: #1d1e26; text-align: center;">Electricity Bill - %s</h2>
				
				<table style="width: 100%%; border-collapse: collapse; margin-top: 20px;">
					<tr>
						<th style="text-align: left; padding: 10px; background: #1d1e26; color: white;">Bill Details</th>
						<th style="text-align: right; padding: 10px; background: #1d1e26; color: white;"></th>
					</tr>
					<tr><td style="padding: 10px;">Bill ID:</td><td style="padding: 10px; text-align: right;">%d</td></tr>
					<tr><td style="padding: 10px;">Billing Date:</td><td style="padding: 10px; text-align: right;">%s</td></tr>
					<tr><td style="padding: 10px;">Owner:</td><td style="padding: 10px; text-align: right;">%s</td></tr>
					<tr><td style="padding: 10px;">Total Consumption:</td><td style="padding: 10px; text-align: right;">%.2f kWh</td></tr>
				</table>

				<h3 style="color: #1d1e26; text-align: center; margin-top: 20px;">Pricing Breakdown</h3>

				<table style="width: 100%%; border-collapse: collapse; margin-top: 10px;">
					<tr>
						<th style="text-align: left; padding: 10px; background: #4d596a; color: white;">Pricelist Details</th>
						<th style="text-align: right; padding: 10px; background: #4d596a; color: white;"></th>
					</tr>
					<tr><td style="padding: 10px;">Pricelist Valid From:</td><td style="padding: 10px; text-align: right;">%s</td></tr>
					<tr><td style="padding: 10px;">Green Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.4f/kWh</td></tr>
					<tr><td style="padding: 10px;">Blue Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.4f/kWh</td></tr>
					<tr><td style="padding: 10px;">Red Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.4f/kWh</td></tr>
					<tr><td style="padding: 10px;">Billing Power Cost:</td><td style="padding: 10px; text-align: right;">$%.4f</td></tr>
					<tr><td style="padding: 10px;">Tax Rate:</td><td style="padding: 10px; text-align: right;">%.2f%%</td></tr>
				</table>

				<h3 style="color: #1d1e26; text-align: center; margin-top: 20px;">Total Amount</h3>

				<table style="width: 100%%; border-collapse: collapse; margin-top: 10px;">
					<tr>
						<th style="text-align: left; padding: 10px; background: #1d1e26; color: white;">Total Price (incl. tax)</th>
						<th style="text-align: right; padding: 10px; background: #1d1e26; color: white;">$%.2f</th>
					</tr>
				</table>

				<div style="text-align: center; margin-top: 30px;">
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Pay Now</a>
				</div>

				<p style="font-size: 14px; color: #888; margin-top: 20px; text-align: center;">If you have already paid, please ignore this message.</p>
			</div>
		</body>
	</html>
	`, billDate, billID, billDate, ownerName, spentPower, validFrom, greenZone, blueZone, redZone, billingPower, tax, price, paymentLink)
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

func GenerateFailedActivationEmailBody(loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">This email has active account already</h2>
					<p style="font-size: 16px; color: #555;">Failed to activate account again!</p>
					<p style="font-size: 16px; color: #555;">Click the button below to go to the login page:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Go to Login</a>
				</div>
			</body>
		</html>
	`, loginLink)
}

func GeneratePropertyApprovalEmailBody(propertyName string, loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Property Approved</h2>
					<p style="font-size: 16px; color: #555;">Congratulations! Your property "<strong>%s</strong>" has been approved.</p>
					<p style="font-size: 16px; color: #555;">Click the button below to log in and manage your property:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Log In</a>
				</div>
			</body>
		</html>
	`, propertyName, loginLink)
}

func GeneratePropertyDeclineEmailBody(propertyName string, reason string, loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Property Declined</h2>
					<p style="font-size: 16px; color: #555;">We regret to inform you that your property "<strong>%s</strong>" has been declined.</p>
					<p style="font-size: 16px; color: #555;">Reason: <strong>%s</strong></p>
					<p style="font-size: 16px; color: #555;">Click the button below to log in:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Log In</a>
				</div>
			</body>
		</html>
	`, propertyName, reason, loginLink)
}

func GenerateOwnershipApprovalEmailBody(householdName string, loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Ownership Approved</h2>
					<p style="font-size: 16px; color: #555;">Congratulations! Your managed to prove that you are the owner of household at "<strong>%s</strong>" .</p>
					<p style="font-size: 16px; color: #555;">Click the button below to log in and see your household:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Log In</a>
				</div>
			</body>
		</html>
	`, householdName, loginLink)
}

func GenerateOwnershipDenialEmailBody(householdName string, reason string, loginLink string) string {
	return fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Ownership Declined</h2>
					<p style="font-size: 16px; color: #555;">Sorry!We regret to inform you that your proof that you are the owner of household at "<strong>%s</strong>" in not valid.</p>
					<p style="font-size: 16px; color: #555;">Reason: <strong>%s</strong></p>
					<p style="font-size: 16px; color: #555;">Click the button below to log in:</p>
					<a href="%s" style="display: inline-block; padding: 12px 24px; background-color: #1d1e26; color: white; text-decoration: none; border-radius: 5px; font-weight: bold;">Log In</a>
				</div>
			</body>
		</html>
	`, householdName, reason, loginLink)
}
