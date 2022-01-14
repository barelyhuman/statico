package main

import (
	"os"
	"strconv"
	"strings"
)

func changeFileExtension(fileName string, checkFor string, replaceWith string) string {
	return strings.Replace(fileName, checkFor, replaceWith, 1)
}

func isMarkdownFile(file os.FileInfo) bool {
	extension := strings.SplitN(file.Name(), ".", -1)
	return extension[len(extension)-1] == "md"
}

func isHTMLFile(file os.FileInfo) bool {
	extension := strings.SplitN(file.Name(), ".", -1)
	return extension[len(extension)-1] == "html"
}

func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}
