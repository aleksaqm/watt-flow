package util

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
	"watt-flow/config"
	"watt-flow/model"

	"github.com/skip2/go-qrcode"
	gomail "gopkg.in/mail.v2"
)

type EmailSender struct {
	emailSecret string
}

func NewEmailSender(env *config.Environment) *EmailSender {
	return &EmailSender{
		emailSecret: env.EmailSecret,
	}
}

func (sender *EmailSender) SendEmail(receiver string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "wattflow12@gmail.com")
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", sender.emailSecret)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	} else {
		return nil
	}
}

func (sender *EmailSender) SendEmailWithQRCode(receiver, subject, body string, qrCodeBytes []byte) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "wattflow12@gmail.com")
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)

	cid := "qr-code"
	message.Embed("qr-code.png", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(qrCodeBytes)
		return err
	}), gomail.SetHeader(map[string][]string{
		"Content-ID": {fmt.Sprintf("<%s>", cid)},
	}))

	message.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", sender.emailSecret)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func GenerateQRCode(link string) ([]byte, error) {
	qr, err := qrcode.New(link, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	buf := new(bytes.Buffer)
	err = qr.Write(256, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write QR code: %w", err)
	}

	return buf.Bytes(), nil
}

func (sender *EmailSender) SendPaymentConfirmation(userEmail, userName string, bill model.Bill) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "wattflow12@gmail.comm")
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", fmt.Sprintf("Payment confirmation #%d", bill.ID))

	body, _ := GenerateSLipBody(bill)
	m.SetBody("text/html", body)

	fileName := fmt.Sprintf("Slip-%d.pdf", bill.ID)
	htmlData, err := GeneratePaymentSlip(&bill, userEmail)
	if err != nil {
		return err
	}

	pdfData, err := GeneratePdfFromHtml(htmlData)
	if err != nil {
		return err
	}
	m.Embed(fileName, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfData)
		return err
	}), gomail.SetHeader(map[string][]string{
		"Content-ID": {fmt.Sprintf("<%s>", "pdf")},
	}))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", sender.emailSecret)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil

}

func GenerateSLipBody(bill model.Bill) (string, error) {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Payment Confirmation</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f7f6;">
  <table border="0" cellpadding="0" cellspacing="0" width="100%%">
    <tr>
      <td style="padding: 20px 0;">
        <table align="center" border="0" cellpadding="0" cellspacing="0" width="600" style="border-collapse: collapse; background-color: #ffffff; border: 1px solid #cccccc;">
          
          <tr>
            <td align="center" style="padding: 40px 0 30px 0; background-color: #004d99; color: #ffffff;">
              <h1 style="font-size: 24px; margin: 0;">Payment Successful!</h1>
            </td>
          </tr>

          <tr>
            <td style="padding: 40px 30px;">
              <h2 style="font-size: 20px; margin: 0 0 20px 0; color: #333333;">Dear,</h2>
              <p style="margin: 0 0 15px 0; font-size: 16px; line-height: 1.5; color: #555555;">
                We hereby confirm that your payment for invoice <strong>#%d</strong> has been successfully processed.
              </p>
              <p style="margin: 0; font-size: 16px; line-height: 1.5; color: #555555;">
                Attached to this email you will find a PDF document with payment details and a completed payment slip.
              </p>
              
              <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="margin-top: 30px; border-top: 1px solid #eeeeee;">
                <tr>
                  <td style="padding: 15px 0; font-size: 16px; color: #555555; width: 50%%;">Payment Amount:</td>
                  <td style="padding: 15px 0; font-size: 16px; color: #111111; font-weight: bold; text-align: right;">%.2f USD</td>
                </tr>
                <tr>
                  <td style="padding: 15px 0; border-top: 1px solid #eeeeee; font-size: 16px; color: #555555;">Billing Period:</td>
                  <td style="padding: 15px 0; border-top: 1px solid #eeeeee; font-size: 16px; color: #111111; text-align: right;">%s</td>
                </tr>
              </table>

            </td>
          </tr>

          <tr>
            <td align="center" style="padding: 30px; background-color: #eeeeee; color: #888888;">
              <p style="margin: 0; font-size: 12px;">
                This is an automatically generated message. Please do not reply to this email.<br>
                © %d Watt-Flow. All rights reserved.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>
`, bill.ID, bill.Price, bill.BillingDate, time.Now().Year()), nil
}

func GeneratePaymentSlip(bill *model.Bill, userEmail string) (string, error) {
	return fmt.Sprintf(
		`
	<html lang="sr">
	<head>
	<meta charset="UTF-8">
	<title>Nalog za uplatu</title>
    <style>
        html, body {
            height: 100%%;
            margin: 0;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: #f9f9f9;
            font-family: 'DejaVu Sans', sans-serif;
        }
    </style
	</head>
	<body>
	<div class="slip">

	<div style="width:800px; background:#fff; border:1px solid #000; padding:20px; font-family:Arial, sans-serif; font-size:14px; color:#000;">

	<div style="text-align:right; font-weight:bold; margin-bottom:10px;">НАЛОГ ЗА УПЛАТУ</div>

	<div style="display:flex;">

	<div style="flex:1; border-right:1px solid #000; padding-right:10px;">
	<div style="margin-bottom:10px;">
	<div style="font-size:12px;">уплатилац</div>
	<div style="border:1px solid #000; height:40px;">%s</div>
	</div>
	<div style="margin-bottom:10px;">
	<div style="font-size:12px;">сврха уплате</div>
	<div style="border:1px solid #000; height:40px;">Рачун за струју %s </div>
	</div>
	<div>
	<div style="font-size:12px;">прималац</div>
	<div style="border:1px solid #000; height:40px;">WattFlow</div>
	</div>

	<div style="margin-top:40px; font-size:12px; border-top:1px solid #000; width:70%%; padding-top:2px;">
		печат и потпис уплатиоца
	</div>
	</div>

	<div style="flex:1; padding-left:10px;">
	<div style="display:flex; gap:10px; margin-bottom:10px; align-items:flex-end;">
	<div style="flex:0.8;">
	<div style="font-size:12px; line-height:14px;">шифра<br>плаћања</div>
	<div style="border:1px solid #000; height:30px; margin-top:2px;">289</div>
	</div>
	<div style="flex:0.8;">
	<div style="font-size:12px;">валута</div>
	<div style="border:1px solid #000; height:30px;">USD</div>
	</div>
	<div style="flex:2;">
	<div style="font-size:12px;">износ</div>
	<div style="border:1px solid #000; height:30px;">%.2f</div>
	</div>
	</div>

	<div style="margin-bottom:10px;">
	<div style="font-size:12px;">рачун примаоца</div>
	<div style="border:1px solid #000; height:30px;">840-1234567890123-45</div>
	</div>

	<div style="display:flex; gap:10px;">
	<div style="flex:0.8;">
	<div style="font-size:12px;">број модела</div>
	<div style="border:1px solid #000; height:30px;">103</div>
	</div>
	<div style="flex:2;">
	<div style="font-size:12px;">позив на број (одобрење)</div>
	<div style="border:1px solid #000; height:30px;">123-4567890123-45</div>
	</div>
	</div>
	</div>
	</div>

	<div style="display:flex; justify-content:space-between; margin-top:40px; font-size:12px;">
	<div style="width:35%%; border-top:1px solid #000; text-align:center; padding-top:2px;">
		место и датум пријема
	</div>
	<div style="width:35%%; border-top:1px solid #000; text-align:center; padding-top:2px;">
		датум валуте
	</div>
	</div>

	</div>
	</div>
	</body>
	</html>`, userEmail, bill.BillingDate, bill.Price,
	), nil
}

func GenerateMonthlyBillEmail(bill *model.Bill) (string, []byte, error) {
	date := time.Time(bill.Pricelist.ValidFrom).Format("2006-01")
	paymentLink := fmt.Sprintf("http://localhost:80/api/pay/%s", strconv.FormatUint(bill.ID, 10))
	qrCodeBytes, err := GenerateQRCode(paymentLink)
	if err != nil {
		return "", nil, fmt.Errorf("error generating QR code: %w", err)
	}
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
					<tr><td style="padding: 10px;">Green Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.3f/kWh</td></tr>
					<tr><td style="padding: 10px;">Blue Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.3f/kWh</td></tr>
					<tr><td style="padding: 10px;">Red Zone Rate:</td><td style="padding: 10px; text-align: right;">$%.3f/kWh</td></tr>
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
        <h3 style="color: #1d1e26; text-align: center; margin-top: 20px;">Scan QR Code to Pay</h3>
            <div style="text-align: center;">
                <img src="cid:qr-code" alt="QR Code for Payment" style="width: 200px; height: 200px; margin-top: 10px;">
            </div>
				<p style="font-size: 14px; color: #888; margin-top: 20px; text-align: center;">If you have already paid, please ignore this message.</p>
			</div>
		</body>
	</html>
	`, date, bill.ID, bill.BillingDate, bill.Owner.Username, bill.SpentPower, date, bill.Pricelist.GreenZone, bill.Pricelist.BlueZone, bill.Pricelist.RedZone, bill.Pricelist.BillingPower, bill.Pricelist.Tax, bill.Price, paymentLink), qrCodeBytes, nil
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
