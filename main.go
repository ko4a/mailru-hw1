package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

func fileSize(file fs.FileInfo) string {
	if file.IsDir() {
		return ""
	}

	fileSize := file.Size()

	if fileSize == 0 {
		return " (empty)"
	}

	return " (" + strconv.FormatInt(fileSize, 10) + "b)"
}

func filterFiles(files []fs.FileInfo, saveFiles bool) []fs.FileInfo {
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	if !saveFiles {
		var res []fs.FileInfo

		for _, file := range files {
			if file.IsDir() {
				res = append(res, file)
			}
		}

		return res

	}

	return files
}

func printFile(out io.Writer, file fs.FileInfo, tabCount int, isLast bool, previousLevelPrefix string) {
	var prefix string

	if isLast {
		prefix = previousLevelPrefix + "└───"
	} else {
		prefix = previousLevelPrefix + "├───"
	}

	fmt.Fprintln(out, prefix+file.Name()+fileSize(file))
}

func printDir(out io.Writer, path string, printFiles bool, level int, prefix string) error {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	files = filterFiles(files, printFiles)

	for idx, file := range files {
		isLast := false

		if idx == len(files)-1 {
			isLast = true
		}

		printFile(out, file, level, isLast, prefix)

		currentPrefix := prefix

		if file.IsDir() {

			if idx != len(files)-1 {
				currentPrefix = currentPrefix + "│\t"
			} else {
				currentPrefix = currentPrefix + "\t"
			}

			filePath := path + string(os.PathSeparator) + file.Name()
			level++
			printDir(out, filePath, printFiles, level, currentPrefix)
			level--
		}

	}

	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return printDir(out, path, printFiles, 0, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
