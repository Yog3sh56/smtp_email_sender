package internal

import (
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type SMTPConfig struct {
	Host        string
	Port        string
	Username    string
	AppPassword string
}

func LoadSMTPConfig() (*SMTPConfig, error) {
	LoadEnv()
	cfg := &SMTPConfig{
		Host:        os.Getenv("SMTP_HOST"),
		Port:        os.Getenv("SMTP_PORT"),
		Username:    os.Getenv("SENDER_EMAIL"),
		AppPassword: os.Getenv("APP_PASSWORD"),
	}

	if cfg.Host == "" || cfg.Port == "" || cfg.Username == "" || cfg.AppPassword == "" {
		return nil, fmt.Errorf("missing SMTP configuration")
	}
	return cfg, nil
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found; relying on system environment")
	}
}

func sendSingleEmailSimplified(cfg *SMTPConfig, to, subject, body string) error {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.AppPassword, cfg.Host)
	msg := buildMessage(cfg.Username, to, subject, body)
	return smtp.SendMail(serverAddr, auth, cfg.Username, []string{to}, []byte(msg))
}

func sendSingleEmailControlled(cfg *SMTPConfig, to, subject, body string) error {
	serverAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	client, err := smtp.Dial(serverAddr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer client.Close()

	tlsConfig := &tls.Config{ServerName: cfg.Host}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls failed: %w", err)
		}
	} else {
		return fmt.Errorf("server %s does not support STARTTLS", serverAddr)
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.AppPassword, cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	if err := client.Mail(cfg.Username); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}
	defer wc.Close()

	msg := buildMessage(cfg.Username, to, subject, body)
	if _, err = wc.Write([]byte(msg)); err != nil {
		return fmt.Errorf("write failed: %w", err)
	}

	if err := client.Quit(); err != nil {
		log.Printf("warning: QUIT returned error for %s: %v", to, err)
	}

	return nil
}

func SendEmail(recipients []string, subject, body string) error {
	cfg, err := LoadSMTPConfig()
	if err != nil {
		return err
	}

	for _, r := range recipients {
		if err := sendSingleEmailSimplified(cfg, r, subject, body); err != nil {
			log.Printf("❌ Failed to send to %s: %v", r, err)
		} else {
			log.Printf("✅ Sent email to %s", r)
		}
	}
	return nil
}

func buildMessage(from, to, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"", // blank line before body
	}
	return strings.Join(headers, "\r\n") + body + "\r\n"
}
