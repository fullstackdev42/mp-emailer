package templates

import (
	"bytes"
	"strings"
	"testing"
)

func TestIndexTemplate(t *testing.T) {
	var buf bytes.Buffer
	data := struct {
		User interface{}
	}{
		User: "TestUser",
	}

	err := IndexTemplate.Execute(&buf, data)
	if err != nil {
		t.Fatalf("Error executing IndexTemplate: %v", err)
	}

	output := buf.String()
	expectedContents := []string{
		"<title>MP Emailer - Contact Your Representative</title>",
		"<h1>MP Emailer</h1>",
		"<h2>Email Your MP: We need REAL AI regulation!</h2>",
		"<h2>Tell Your MP: Keep the Internet Open and Free</h2>",
		"<h2>About MP Emailer</h2>",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("IndexTemplate output does not contain expected content: %s", expected)
		}
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
	expectedContents := []string{
		"<title>MP Emailer - Contact Your Representative</title>",
		"<h1>Email to Member of Parliament</h1>",
		"<p>To: mp@example.com</p>",
		"<pre>This is a test email content.</pre>",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("EmailTemplate output does not contain expected content: %s", expected)
		}
	}
}

func TestLoginTemplate(t *testing.T) {
	var buf bytes.Buffer
	err := LoginTemplate.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("Error executing LoginTemplate: %v", err)
	}

	output := buf.String()
	expectedContents := []string{
		"<title>MP Emailer - Contact Your Representative</title>",
		"<h1>Login</h1>",
		`<form action="/login" method="post">`,
		`<input type="text" id="username" name="username" required>`,
		`<input type="password" id="password" name="password" required>`,
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("LoginTemplate output does not contain expected content: %s", expected)
		}
	}
}
