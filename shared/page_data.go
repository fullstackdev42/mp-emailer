package shared

// PageData represents the common data structure for page rendering
type PageData struct {
	Title    string
	Content  interface{}
	Error    string
	Messages []string
}
