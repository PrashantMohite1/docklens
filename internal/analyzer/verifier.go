package analyzer

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func Verify_file_sha256_in_container(imageName string, file string) {

	pairs := strings.Split(file, ",")

	for _, pair := range pairs {
		spl2 := strings.Split(pair, ":")
		localpath := spl2[0]
		imgpath := spl2[1]
		Verify_ech_file_sha256_in_container(imageName, localpath, imgpath)
	}

}

func Verify_ech_file_sha256_in_container(imageName string, filePath string, imagepath string) {

	start := time.Now()
	localhash := Get_local_files_sha256(filePath)

	containerhash, err := Run_in_container(imageName, imagepath)
	if err != nil {
		log.Fatal(err)
	}

	localHex := fmt.Sprintf("%x", localhash)
	containerHex := strings.TrimSpace(string(containerhash))

	if localHex == containerHex {
		fmt.Printf("%s : SHA256 hashes match!\n", filePath)
	} else {
		fmt.Printf("%s : SHA256 hashes do not match!\n", filePath)
	}

	elapsed := time.Since(start)
	fmt.Printf("Verification took %s\n", elapsed)

}
