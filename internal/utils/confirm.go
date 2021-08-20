package utils

import (
	"fmt"
	"log"
  "os"
)

func UserConfirmation() {
	var response string
  fmt.Print("Do you want to proceed [y/n] ? ")

	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}

	if response == "y" ||  response == "Y" {
		return
	} else {
		os.Exit(1)
	}
}
