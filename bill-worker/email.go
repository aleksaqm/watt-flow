package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/skip2/go-qrcode"
	gomail "gopkg.in/mail.v2"
)

type EmailSender struct {
	emailSecret string
}

func NewEmailSender(secret string) *EmailSender {
	return &EmailSender{
		emailSecret: secret,
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

func GenerateMonthlyBillEmail(bill Bill, pricelist Pricelist, user string) (string, []byte, error) {
	date := time.Time(pricelist.ValidFrom).Format("2006-01")
	paymentLink := fmt.Sprintf("http://localhost:80/api/pay/%s", "test")
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
	`, date, bill.BillingDate, user, bill.SpentPower, date, pricelist.GreenZone, pricelist.BlueZone, pricelist.RedZone, pricelist.BillingPower, pricelist.Tax, bill.Price, paymentLink), qrCodeBytes, nil
}
