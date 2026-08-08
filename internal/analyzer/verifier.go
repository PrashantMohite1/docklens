package analyzer

import (
	"fmt"
	"log/slog"
	"strings"
)

func Verify_dir_in_container(imgname string, dirpath string) {

	multi_dirpaths := strings.Split(dirpath, ",")

	for _, directory := range multi_dirpaths {

		idx := strings.LastIndex(directory, ":")
		if idx == -1 {
			slog.Error("invalid dir mapping", "directory", directory)
			return
		}
		localdir := directory[:idx]
		imagedir := directory[idx+1:]
		slog.Info("Directory mapping", "local_dir", localdir, "image_dir", imagedir)

		cmd := []string{
			"/bin/sh",
			"-c",
			"find " + imagedir + " -type f -exec sha256sum {} + | awk '{print $1}' | sort | sha256sum",
		}

		localdir_hash := Checkdir_hash(localdir)

		containerdir_hash, err := Run_in_container(imgname, imagedir, cmd)
		if err != nil {
			slog.Error("Error while running cmd in container", "error", err)
		}

		localHex := fmt.Sprintf("%x", localdir_hash)
		containerHex := strings.TrimSpace(string(containerdir_hash))

		if localHex == containerHex {
			slog.Info("Directory hashes match!")
		} else {
			slog.Info("Directory hashes do not match!")
		}

	}
}

func Verify_file_sha256_in_container(imageName string, file string) {

	pairs := strings.Split(file, ",")

	for _, pair := range pairs {
		idx := strings.LastIndex(pair, ":")
		if idx == -1 {
			slog.Error("invalid file mapping", "filemap", pair)
			return
		}
		localpath := pair[:idx]
		imagepath := pair[idx+1:]
		Verify_ech_file_sha256_in_container(imageName, localpath, imagepath)
	}

}

func Verify_ech_file_sha256_in_container(imageName string, filePath string, imagepath string) {

	localhash := Get_local_files_sha256(filePath)

	cmd := []string{"/bin/sh", "-c", "sha256sum " + imagepath}

	containerhash, err := Run_in_container(imageName, imagepath, cmd)
	if err != nil {
		slog.Error("Error while running command in container", err)
	}

	localHex := fmt.Sprintf("%x", localhash)
	containerHex := strings.TrimSpace(string(containerhash))

	if localHex == containerHex {
		slog.Info("SHA256 hashes match", "file", filePath)
	} else {
		slog.Error("SHA256 hashes do not match", "file", filePath)
	}

}
