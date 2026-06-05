package main

import (
	"fmt"

	"github.com/9thanaphat/go/thanaphat"
	"github.com/google/uuid"
)

func sayHello() {
	fmt.Println("Hello, World!")
}

func main() {
	sayHello()
	thanaphat.SayHelloThanaphat()
	id := uuid.New()
	fmt.Printf("Generated UUID: %s\n", id.String())
}
