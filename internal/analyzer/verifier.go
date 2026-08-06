package analyzer

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func Verify_dir_in_container(imgname string, dirpath string) {

	idx := strings.LastIndex(dirpath, ":")
	if idx == -1 {
		fmt.Println("invalid dir mapping")
		return
	}
	localdir := dirpath[:idx]
	imagedir := dirpath[idx+1:]
	fmt.Printf("local dir : %s  image dir : %s \n", localdir, imagedir)

	cmd := []string{
		"/bin/sh",
		"-c",
		"find " + imagedir + " -type f -exec sha256sum {} + | awk '{print $1}' | sort | sha256sum",
	}

	localdir_hash := Checkdir_hash(localdir)

	containerdir_hash, err := Run_in_container(imgname, imagedir, cmd)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Local dir hash : %x \n", localdir_hash)
	fmt.Printf("Container dir hash : %x \n", containerdir_hash)

}

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

	cmd := []string{"/bin/sh", "-c", "sha256sum " + imagepath}

	containerhash, err := Run_in_container(imageName, imagepath, cmd)
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
