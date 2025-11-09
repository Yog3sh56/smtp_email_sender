package main

import (
	email "Email_SMTP_App/src/cmd/internal"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	recipients := []string{
		"test1@gmail.com",
		"test2@hotmail.co.uk",
		"test3@gmail.com",
	}

	subject := "Sample Subject"
	body := `
Sample Email Body
`

	// Optionally override recipients using TEST_RECIPIENTS env var
	if envList := os.Getenv("TEST_RECIPIENTS"); envList != "" {
		recipients = splitCSV(envList)
	}

	fmt.Printf("Sending email to %d recipients...\n", len(recipients))

	if err := email.SendEmail(recipients, subject, body); err != nil {
		log.Fatalf("Failed to send emails: %v", err)
	}

	fmt.Println("🎉 All emails processed.")
}

func splitCSV(s string) []string {
	var out []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
