package main

import (
	"fmt"
)

func (c *Connection) PassThru(message string) {
	fmt.Printf("[Debug] Message received: ‘%s’.\n", message)
}