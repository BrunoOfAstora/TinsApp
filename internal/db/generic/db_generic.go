package generic

import (
	"log"
	"os"
	"path/filepath"
)

func DbFilePath() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	dbFolder := filepath.Join(dir, "TinsAppDB")
	err = os.MkdirAll(dbFolder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	dbPath := filepath.Join(dbFolder, "tins_db.db")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}
	return dbPath
}
