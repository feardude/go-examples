package main

import "log"

func main() {
	log.Println("Started FX service")
	InitDB()

	defer ShutdownDB()

	log.Println("Finished FX service")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
