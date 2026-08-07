package analyzer

import (
	"fmt"
	"log"
	"strings"
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

	localHex := fmt.Sprintf("%x", localdir_hash)
	containerHex := strings.TrimSpace(string(containerdir_hash))

	fmt.Printf("Local dir hash : %s\n", localHex)
	fmt.Printf("Container dir hash : %s\n", containerHex)

	if localHex == containerHex {
		fmt.Println("Directory hashes match!")
	} else {
		fmt.Println("Directory hashes do not match!")
	}

}

func Verify_file_sha256_in_container(imageName string, file string) {

	pairs := strings.Split(file, ",")

	// for _, pair := range pairs {
	// 	fmt.Println(pair)
	// }

	for _, pair := range pairs {
		idx := strings.LastIndex(pair, ":")
		if idx == -1 {
			fmt.Println("invalid dir mapping")
			return
		}
		localpath := pair[:idx]
		imagepath := pair[idx+1:]

		// fmt.Printf("Local path : %s \n\n Image path : %s \n", localpath, imagepath)

		Verify_ech_file_sha256_in_container(imageName, localpath, imagepath)
	}

}

func Verify_ech_file_sha256_in_container(imageName string, filePath string, imagepath string) {

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

}
