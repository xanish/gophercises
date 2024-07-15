package main

import (
	"fmt"
	"github.com/xanish/gophercises/strings_and_bytes"
)

func main() {
	fmt.Println(strings_and_bytes.Encrypt("There's-a-starman-waiting-in-the-sky", 3))
	fmt.Println(strings_and_bytes.CountWords("helloWorld"))
}
