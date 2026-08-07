package analyzer

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
)

func Get_local_files_sha256(localfilepath string) []byte {

	file, err := os.Open(localfilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	return hash.Sum(nil)

}
