package helpers

import (
	"encoding/gob"
)

// Types use by gob as interface{} have to be registered
func init() {
	gob.Register(Person{})
}

/* Structs for the messages: */

type Person struct {
	Name   string
	Adress string
	Mail   string
}
