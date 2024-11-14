package util

import (
	"encoding/base64"
	"fmt"
	"os"
)

func SaveFile(filename string, data string, filetype string, fileFolder string) (string, error) {
	fileData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 image: %w", err)
	}
	fileName := filename + "." + filetype
	filePath := "/app/data/" + fileFolder
	_ = os.Mkdir(filePath, os.ModePerm)

	filePath = filePath + "/" + fileName
	file, err3 := os.Create(filePath)
	if err3 != nil {
		return "", fmt.Errorf("failed to create file for profile images: %w", err3)
	}
	_, err4 := file.Write(fileData)
	if err4 != nil {
		return "", fmt.Errorf("failed to write image to file: %w", err4)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	return filePath, nil
}
