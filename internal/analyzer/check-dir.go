package analyzer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func hash_check_in_files(filepath string) []byte {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		fmt.Println(err)
	}
	return hash.Sum(nil)

}

func Checkdir_hash(dirpath string) [32]byte {

	var hashes []string

	filepath.WalkDir(dirpath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		h := hash_check_in_files(path)
		hashes = append(hashes, fmt.Sprintf("%x", h))

		return nil
	})

	sort.Strings(hashes)

	manifest := strings.Join(hashes, "\n")
	if manifest != "" {
		manifest += "\n"
	}

	final := sha256.Sum256([]byte(manifest))
	return final

}
