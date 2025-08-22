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

	if os.MkdirAll(test_directory_path, os.ModePerm) != nil {
		fmt.Print(`Error automatically creating the directory. Update the environemnt"
variable 'gotestdir' and/or make the directory yourself.`)

		return
	}

	// ===========================================================================
	// Next step: parse command line input to figure out what the user wants to
	// do, such as make a new test, list a test, delete a test, or whatever else.
	// I'm thinking that my decision to get rid of bubbletea was ill-advised
	// because the first thing that came to mind was a way to search all the tests
	// but I can come back to making a cli if and when the basic functionality is
	// done.
	// ===========================================================================
	fmt.Println("So far, so good!")
}
