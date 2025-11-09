package internal

import (
	"os"
	"strings"
	"testing"
)

// ✅ Unit test for successful SMTP config load
func TestLoadSMTPConfig_Success(t *testing.T) {
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SENDER_EMAIL", "test@example.com")
	os.Setenv("APP_PASSWORD", "dummy-password")

	cfg, err := LoadSMTPConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Host != "smtp.gmail.com" {
		t.Errorf("expected host smtp.gmail.com, got %s", cfg.Host)
	}
	if cfg.Port != "587" {
		t.Errorf("expected port 587, got %s", cfg.Port)
	}
	if cfg.Username != "test@example.com" {
		t.Errorf("expected username test@example.com, got %s", cfg.Username)
	}
	if cfg.AppPassword != "dummy-password" {
		t.Errorf("expected app password dummy-password, got %s", cfg.AppPassword)
	}
}

// ✅ Unit test for missing env vars
func TestLoadSMTPConfig_MissingValues(t *testing.T) {
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SENDER_EMAIL")
	os.Unsetenv("APP_PASSWORD")

	_, err := LoadSMTPConfig()
	if err == nil {
		t.Fatal("expected error for missing SMTP config, got nil")
	}
}

// ✅ Unit test for message building
func TestBuildMessage(t *testing.T) {
	from := "sender@example.com"
	to := "receiver@example.com"
	subject := "Test Subject"
	body := "Hello world!"

	msg := buildMessage(from, to, subject, body)

	// Check that required headers and body exist
	tests := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"Content-Type: text/plain",
		body,
	}

	for _, want := range tests {
		if !strings.Contains(msg, want) {
			t.Errorf("expected message to contain %q", want)
		}
	}
}

// ⚙️ Optional integration test for real SMTP send
// Run manually with: RUN_SMTP_TESTS=true go test -v ./internal
func TestSendEmail_Integration(t *testing.T) {
	if os.Getenv("RUN_SMTP_TESTS") != "true" {
		t.Skip("Skipping live SMTP test. Set RUN_SMTP_TESTS=true to enable.")
	}

	recipients := []string{"your.email@gmail.com"}
	subject := "Integration Test Email"
	body := "This is a test email sent from Go."

	err := SendEmail(recipients, subject, body)
	if err != nil {
		t.Fatalf("failed to send email: %v", err)
	}
}
