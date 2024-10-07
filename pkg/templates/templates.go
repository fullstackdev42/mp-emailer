package templates

import "html/template"

var (
	IndexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>MP Emailer</title>
</head>
<body>
    <h1>MP Emailer</h1>
    <form action="/submit" method="post">
        <label for="postalCode">Enter your postal code:</label>
        <input type="text" id="postalCode" name="postalCode" required>
        <input type="submit" value="Find MP">
    </form>
</body>
</html>
`))

	EmailTemplate = template.Must(template.New("email").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>MP Email</title>
</head>
<body>
    <h1>Email to MP</h1>
    <p>To: {{.Email}}</p>
    <pre>{{.Content}}</pre>
</body>
</html>
`))
)
