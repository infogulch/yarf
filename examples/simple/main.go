package main

import (
	"github.com/yarf-framework/yarf"
)

// Define a simple resource
type Hello struct {
	yarf.Resource // Extend the yarf.Resource by composition
}

// Implement the GET handler
func (h *Hello) Get(c *yarf.Context) error {
	c.Render("Hello world!")

	return nil
}

// Entry point of the executable application
// It runs a default server listening on http://localhost:8080
func main() {
	// Create a new empty YARF server
	y := yarf.New()

	// Add route/resource
	y.Add("/", new(Hello))

	// Start server listening on port 8080
	y.Start(":8080")
}
