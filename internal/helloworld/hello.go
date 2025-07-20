package helloworld

import "fmt"

// Greet returns a greeting for the named person.
func Greet(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
