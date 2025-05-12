package dirhandler

import (
	"fmt"
	"log"
	"os"
)

func CheckOutputDir(o string) {

	if _, err := os.Stat(o); os.IsNotExist(err) {
		var answer string

		// Ask user if he wants to create the directory if it does not exist
		for {
			fmt.Println("Directory " + o + "does not exist \n Would you like it to be created ?(yes/no):")
			fmt.Scanln(&answer)
			if answer == "yes" || answer == "no" {
				break
			}
			fmt.Println("Invalid input. Please try again.")

		}

		// if answer yes make the directory otherwise exit
		if answer == "yes" {
			mkdirErr := os.MkdirAll(o, 0755)
			if mkdirErr != nil {
				log.Fatalln("There was a problem with creating directories:", mkdirErr)
			}

		} else {
			os.Exit(1)
		}

	}

}
