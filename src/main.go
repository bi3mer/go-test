package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const TEST_DIRECTORY = "gotestdir"

func main() {
	// ===========================================================================
	// Set up and get user test directory
	// ===========================================================================
	var test_directory_path = os.Getenv(TEST_DIRECTORY)
	if test_directory_path == "" {
		fmt.Println("Test directory must be set (e.g. `export gotestdir=~/tests`).")
		return
	}

	if strings.HasPrefix(test_directory_path, "~/") {
		home, error_get_user_home_dir := os.UserHomeDir()
		if error_get_user_home_dir != nil {
			fmt.Println("Unexpected error getting user home dir: ", error_get_user_home_dir)
			return
		}

		test_directory_path = filepath.Join(home, test_directory_path[2:])
	}

	error_make_test_directory := os.MkdirAll(test_directory_path, os.ModePerm)
	if error_make_test_directory != nil {
		fmt.Print(`Error automatically creating the directory. Update the environemnt"
variable 'gotestdir' and/or make the directory yourself.`)

		return
	}

	// ===========================================================================
	// Se
	// ===========================================================================
	fmt.Println("So far, so good!")
}
