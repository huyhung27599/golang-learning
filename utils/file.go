package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

var allowExts = map[string]bool{
	"jpg": true,
	"jpeg": true,
	"png": true,
	"gif": true,
	"webp": true,
}

const maxFileSize = 5<<20

var allowMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png": true,
	"image/gif": true,
	"image/webp": true,
}

func ValidateAndSaveFile( fileHeader *multipart.FileHeader, uploadDir string) (string, error) {
 ext :=	strings.ToLower(filepath.Ext(fileHeader.Filename))

 if !allowExts[ext] {
	return "", fmt.Errorf("invalid file extension")
 }

 if fileHeader.Size > maxFileSize {
	return "", fmt.Errorf("file size is too large")
 }

 file, err := fileHeader.Open()
 if err != nil {
	return "", fmt.Errorf("failed to open file")
 }
 defer file.Close()

 buffer := make([]byte, 512)
  _, err = file.Read(buffer)
  if err != nil {
	return "", fmt.Errorf("failed to read file")
  }

  mimeType := http.DetectContentType(buffer)

  if !allowMimeTypes[mimeType] {
	return "", fmt.Errorf("invalid file mime type")
  }

  fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

  
  if e := os.MkdirAll(uploadDir, os.ModePerm); e != nil {
	return "", fmt.Errorf("failed to create upload directory")
  }

  savePath := filepath.Join(uploadDir, fileName)
  if err := saveFile(fileHeader, savePath); err != nil {
	return "", fmt.Errorf("failed to save file")
  }

 return fileName,nil
}

func saveFile(fileHeader *multipart.FileHeader, destination string) ( error) {
	 src, err := fileHeader.Open()
	 if err != nil {
		return fmt.Errorf("failed to open file")
	 }
	 defer src.Close()

	 dst, err := os.Create(destination)
	 if err != nil {
		return fmt.Errorf("failed to create file")
	 }
	 defer dst.Close()

	 _, err = io.Copy(dst, src)
	 if err != nil {
		return fmt.Errorf("failed to copy file")
	 }
	 return nil
}