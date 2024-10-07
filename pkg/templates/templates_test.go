package templates

import (
	"bytes"
	"testing"
)

func TestIndexTemplate(t *testing.T) {
	var buf bytes.Buffer
	data := struct{}{} // No data needed for the index template

	err := IndexTemplate.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Error executing IndexTemplate: %v", err)
	}

	output := buf.String()
	if !contains(output, "<title>MP Emailer</title>") {
		t.Errorf("IndexTemplate output does not contain expected title")
	}
	if !contains(output, `<form action="/submit" method="post">`) {
		t.Errorf("IndexTemplate output does not contain expected form")
	}
}

func TestEmailTemplate(t *testing.T) {
	var buf bytes.Buffer
	data := struct {
		Email   string
		Content string
	}{
		Email:   "mp@example.com",
		Content: "This is a test email content.",
	}

	err := EmailTemplate.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Error executing EmailTemplate: %v", err)
	}

	output := buf.String()
	if !contains(output, "<title>Email to Member of Parliament</title>") {
		t.Errorf("EmailTemplate output does not contain expected title")
	}
	if !contains(output, `<p>To: mp@example.com</p>`) {
		t.Errorf("EmailTemplate output does not contain expected email address")
	}
	if !contains(output, `<pre>This is a test email content.</pre>`) {
		t.Errorf("EmailTemplate output does not contain expected email content")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
