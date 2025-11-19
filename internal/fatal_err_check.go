package internal

import (
	"log"
)

func FatalErrChecking(err error) error {
	if err != nil {
		log.Fatal(err)
	}

	return err
}
