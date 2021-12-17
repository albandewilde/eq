package main

import (
	"fmt"
)

func main() {
	hash, _ := getHash("./files/file0")
	fmt.Println(hash)

	duplication, err := Duplicates("/home/alban/Desktop/")
	fmt.Println(err)
	fmt.Println(duplication)
}
