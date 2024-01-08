package main

import (
	"fmt"
	"io"
	"os"

	"github.com/xwb1989/sqlparser"
)

func main() {

	r, _ := os.Open("file.sql")
	defer r.Close()
	tokens := sqlparser.NewTokenizer(r)
	fmt.Print("\n")
	for {
		_, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		fmt.Printf("\rReading: %d", tokens.Position)

	}
	fmt.Print("\n")
}
