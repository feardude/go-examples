package main

import "log"

func main() {
	log.Println("Wololo")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
