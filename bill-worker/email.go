package main

import (
	"bytes"
	"encoding/base64"
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

func (sender *EmailSender) SendEmailWithPDFAttachment(receiver, subject, body string, pdfBytes []byte, filename string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "wattflow12@gmail.com")
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	// Attach PDF
	message.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(pdfBytes)
		return err
	}))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", sender.emailSecret)

	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
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

func GenerateMonthlyBillHTML(bill Bill, pricelist Pricelist, user string, householdCN string) (string, error) {
	date := time.Time(pricelist.ValidFrom).Format("2006-01")
	paymentLink := fmt.Sprintf("http://localhost:80/bills/pay/%s", bill.PaymentReference)

	// Generate QR code
	qrCodeBytes, err := GenerateQRCode(paymentLink)
	if err != nil {
		return "", fmt.Errorf("error generating QR code: %w", err)
	}

	// Convert QR code to base64 for embedding in HTML
	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeBytes)

	htmlContent := fmt.Sprintf(`
	<html>
		<body style="font-family: Arial, sans-serif; background: #f4f4f4; color: #333; padding: 40px; text-align: center;">
			<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: left;">
				<h2 style="color: #1d1e26; text-align: center;">Electricity Bill - %s - household: %s</h2>

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
					<img src="data:image/png;base64,%s" alt="QR Code for Payment" style="width: 200px; height: 200px; margin-top: 10px;">
				</div>

				<p style="font-size: 14px; color: #888; margin-top: 20px; text-align: center;">If you have already paid, please ignore this message.</p>
			</div>
		</body>
	</html>
	`, bill.BillingDate, householdCN, bill.BillingDate, user, bill.SpentPower, date, pricelist.GreenZone, pricelist.BlueZone, pricelist.RedZone, pricelist.BillingPower, pricelist.Tax, bill.Price, paymentLink, qrCodeBase64)

	return htmlContent, nil
}

func (sender *EmailSender) SendMonthlyBillPDF(receiver string, bill Bill, pricelist Pricelist, user string, householdCN string) error {
	htmlContent, err := GenerateMonthlyBillHTML(bill, pricelist, user, householdCN)
	if err != nil {
		return fmt.Errorf("error generating HTML content: %w", err)
	}

	pdfBytes, err := GeneratePdfFromHtml(htmlContent)
	if err != nil {
		return fmt.Errorf("error generating PDF: %w", err)
	}

	filename := fmt.Sprintf("electricity_bill_%s_%s.pdf", householdCN, bill.BillingDate)

	subject := fmt.Sprintf("Electricity Bill - %s - %s", bill.BillingDate, householdCN)
	emailBody := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; background: linear-gradient(90deg, #1d1e26 0%%, #4d596a 100%%); color: #333; padding: 40px; text-align: center;">
				<div style="background: white; max-width: 600px; margin: 0 auto; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
					<h2 style="color: #333; margin-bottom: 20px;">Monthy Bill - %s</h2>
					<p style="font-size: 16px; color: #555;">Dear %s, please find monthly bill for household: %s attached to this email.</p>
				</div>
			</body>
		</html>
`, bill.BillingDate, user, householdCN)

	// Send email with PDF attachment
	return sender.SendEmailWithPDFAttachment(receiver, subject, emailBody, pdfBytes, filename)
}

func GenerateMonthlyBillEmail(bill Bill, pricelist Pricelist, user string, householdCN string) (string, []byte, error) {
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
				<h2 style="color: #1d1e26; text-align: center;">Electricity Bill - %s - household: %s</h2>

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
	`, bill.BillingDate, householdCN, bill.BillingDate, user, bill.SpentPower, date, pricelist.GreenZone, pricelist.BlueZone, pricelist.RedZone, pricelist.BillingPower, pricelist.Tax, bill.Price, paymentLink), qrCodeBytes, nil
}
