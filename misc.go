package main

import (
	"io"
	"log"
)

func closeResource(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}
