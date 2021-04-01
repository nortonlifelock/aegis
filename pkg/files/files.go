package files

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ExecuteThroughDirectory takes a first order function as a parameter and executes it against every single file in a directory
func ExecuteThroughDirectory(directory string, recursive bool, exec func(fpath string, file os.FileInfo) error) (err error) {
	var filesWithinDirectory []os.FileInfo
	var filePaths []string
	var middleError error

	filesWithinDirectory, filePaths, err = collectFilesFromDirectory(directory, recursive)
	if err == nil {
		for fileNum, file := range filesWithinDirectory {

			fpath := filepath.Join(filePaths[fileNum])
			middleError = executeOnFile(fpath, file, exec)
			if middleError != nil { // so we don't overwrite a valid error with nil
				if err == nil {
					err = errors.Errorf("%s: %s", file.Name(), middleError.Error())
				} else {
					err = errors.Errorf("%s: %s\n%s", file.Name(), middleError.Error(), err.Error())
				}
			}
		}
	}
	return err
}

// validPath identifies if the given path is valid given the first order function
func validPath(path string, verify func(finfo os.FileInfo) bool) (valid bool) {
	var err error
	var finfo os.FileInfo

	if len(path) > 0 {
		if finfo, err = os.Stat(path); err == nil {
			if finfo != nil && verify(finfo) {
				valid = true
			}
		}
	}

	return valid
}

// ValidFile identifies if the given path exists and points to a file
func ValidFile(path string) (valid bool) {
	return validPath(path, func(finfo os.FileInfo) bool {
		return !finfo.IsDir()
	})
}

// ValidDir identifies if the given path exists and points to a directory
func ValidDir(path string) (valid bool) {
	return validPath(path, func(finfo os.FileInfo) bool {
		return finfo.IsDir()
	})
}

func executeOnFile(fpath string, file os.FileInfo, exec func(fpath string, file os.FileInfo) error) (err error) {
	// Ignore hidden files
	if len(file.Name()) > 0 && file.Name()[0] != '.' {
		if exec != nil {
			// This path is a file not a directory so execute the
			// function that was passed against this file
			err = exec(fpath, file)
		} else {
			err = errors.New("Method to Execute Through Directory cannot be nil")
		}
	}
	return err
}

func collectFilesFromDirectory(directory string, recursive bool) (fileSlice []os.FileInfo, filePaths []string, err error) {
	if len(directory) > 0 {
		// Get the directory information and ensure the directory is actually a directory
		var finfo os.FileInfo
		if finfo, err = os.Stat(directory); err == nil {
			if finfo != nil && finfo.IsDir() {
				// Read the files from the directory structure
				var listOfFiles []os.FileInfo
				if listOfFiles, err = ioutil.ReadDir(directory); err == nil {

					// Verify that files were returned from the directory read
					if listOfFiles != nil && len(listOfFiles) > 0 {
						for _, file := range listOfFiles {

							filePath := filepath.Join(directory, file.Name())

							if !file.IsDir() { //Is a file
								fileSlice = append(fileSlice, file)
								filePaths = append(filePaths, filePath)
								//This is covered by a separate test case
							} else if recursive { //Is a directory
								var newFiles []os.FileInfo
								var newPaths []string
								newFiles, newPaths, err = collectFilesFromDirectory(filePath, recursive)
								fileSlice = append(fileSlice, newFiles...)
								filePaths = append(filePaths, newPaths...)
							}
						}
					}
				}
			} else {
				err = errors.Errorf("Path used is not a valid file or directory. PATH: %s", directory)
			}
		}
	} else {
		err = errors.New("Path cannot be empty")
	}

	return fileSlice, filePaths, err
}

// GetFileContents returns a bytes buffer holding the contents of the file located at the path
func GetFileContents(path string) (buff *bytes.Buffer, err error) {
	buff = bytes.NewBuffer(nil)

	var f *os.File
	if f, err = os.Open(path); err == nil {

		// Close the file when this method is finished
		defer func() {
			_ = f.Close()
		}()

		// Copy the buffer array from the file and store it in buff
		_, err = io.Copy(buff, f)
	}

	return buff, err
}

// GetStringFromFile returns a string holding the contents of the file located at the path
func GetStringFromFile(path string) (contents string, err error) {
	if len(path) > 0 {
		var file *bytes.Buffer
		if file, err = GetFileContents(path); err == nil {

			// Convert the bytes from the file to a string
			contents = string(file.Bytes())
		}
	} else {
		err = errors.New("Invalid File Path")
	}

	return contents, err
}

// WriteFile writes the contents of the contents parameter to the file located at the end of the path
func WriteFile(path string, contents string) (err error) {
	err = ioutil.WriteFile(path, []byte(contents), 0664)
	return err
}
