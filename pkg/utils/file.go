package utils

import (
	"Open_IM/pkg/common/constant"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

// Determine whether the given path is a folder
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Determine whether the given path is a file
func IsFile(path string) bool {
	return !IsDir(path)
}

// Create a directory
func MkDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GetNewFileNameAndContentType(fileName string, fileType int) (string, string) {
	suffix := path.Ext(fileName)
	newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), fileName)
	contentType := ""
	if fileType == constant.ImageType {
		contentType = "image/" + suffix[1:]
	}
	return newName, contentType
}

func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)
	switch {
	case bytes >= EXABYTE:
		unit = "E"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "P"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}
	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// LoadByteFiles - read files from folder and return bytes, filtered by extension
func LoadByteFiles(dirname string, ext string) ([]byte, error) {
	var strCode []byte

	err := filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		// ext = .js
		if !f.IsDir() && strings.HasSuffix(f.Name(), ext) {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			strCode = append(strCode, b...)
		}
		return nil
	})
	return strCode, err
}

// LoadListFiles - loads a directory recursivily and returns a list
// path : directory path to start create the list
// ext : extension to filter, only files with a especific extension will be included into list , other files will bo ignored
// removeExtension : if true remove extension form the file name
// Example : LoadListFiles("/Users/test", ".html", true)
func LoadListFiles(path string, ext string, removeExtension bool) ([]string, error) {
	var list []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		if f.IsDir() {
			listDir, err := LoadListFiles(path+"/"+f.Name(), ext, removeExtension)
			if err != nil {
				return list, err
			}

			for _, s := range listDir {

				if removeExtension == true {
					list = append(list, f.Name()+"/"+strings.Replace(s, ext, "", 1))
				} else {
					list = append(list, f.Name()+"/"+s)
				}

			}

		} else {
			fileExt := filepath.Ext(f.Name())
			if fileExt == ext {
				if removeExtension == true {
					list = append(list, strings.Replace(f.Name(), ext, "", 1))
				} else {
					list = append(list, f.Name())
				}
			}
		}

	}

	return list, err
}

// LoadFilesInfo - loads a directory recursivily and returns a map[string]interface - list of infos
// path : directory path to start create the list
// Example : LoadFilesInfo("/Users/test")
func LoadFilesInfo(path string) ([]map[string]interface{}, error) {

	list := []map[string]interface{}{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return list, err
	}

	for _, f := range files {
		var listItem = make(map[string]interface{})
		if f.IsDir() {
			listItem["name"] = f.Name()
			abs := path + string(filepath.Separator) + f.Name()
			listItem["absolutePath"], _ = filepath.Abs(abs)
			listItem["extension"] = ""
			listItem["path"] = abs
			listItem["isDir"] = true
			listDir, err := LoadFilesInfo(path + "/" + f.Name())
			if err != nil {
				return list, err
			}
			listItem["childs"] = listDir

		} else {
			fileExt := filepath.Ext(f.Name())
			listItem["name"] = f.Name()
			abs := path + string(filepath.Separator) + f.Name()
			listItem["absolutePath"] = abs
			listItem["extension"] = fileExt[1:]
			listItem["path"] = abs
			listItem["isDir"] = false
			listItem["childs"] = []map[string]interface{}{}
		}
		list = append(list, listItem)
	}
	return list, err
}

// LoadBytesDir - read files from folder
func LoadBytesDir(dirname string) ([]byte, error) {
	var strCode []byte

	err := filepath.Walk(dirname, func(path string, f os.FileInfo, err error) error {
		// ext = .js
		if !f.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			strCode = append(strCode, b...)
		}
		return nil
	})
	return strCode, err
}

// LoadJson  - Load a json file and return into a inrterface
// Example : LoadJson("./file.json", &obj)
func LoadJson(path string, obj interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &obj)
	if err != nil {
		return err
	}
	return nil
}

// SaveJson  - Convert a interface into json and save a file
func SaveJson(path string, obj interface{}) error {

	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0644)
}
func SaveString(path string, content *[]byte) error {
	return ioutil.WriteFile(path, *content, 0644)
}
func ReadString(path string) (string, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return Bytes2string(file), nil
}

// RemoveDuplicates - remove duplicate strings from slice string
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// GetCWD - return working dir
func GetCWD() string {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR:  Could not get working directory.\n")
		fmt.Fprintf(os.Stderr, "ERROR-MESSAGE:%v\n", err)
		os.Exit(4)
	}
	return currentWorkingDirectory
}

// RenameIfExists - rename a file if exists
func RenameIfExists(path string) {
	os.Rename(path, fmt.Sprintf("%s-Pre-%s", path, GetTimeStamp()))
}

const TIME_LAYOUT = "Jan-02-2006_15-04-05-MST"

// GetTimeStamp - return timeStamp string with current date
func GetTimeStamp() string {
	now := time.Now()
	return now.Format(TIME_LAYOUT)
}
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
