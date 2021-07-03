package utilities

import (
	"fmt"
	"mime/multipart"
    "os"
    "path/filepath"
	"strings"
    "time"
)

func FindInSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func ValidateUploadFile(file *multipart.FileHeader, maxSizeKb int64, mimeTypes []string) Response {
	res := Response{}

	if (file.Size / 1024) > maxSizeKb {
		res.Message = fmt.Sprintf("Tamaño del archivo debe ser menor o igual a %d KB", maxSizeKb)
		return res
	}

	fileExt := strings.Trim(filepath.Ext(file.Filename), ".")
	_, found := FindInSlice(mimeTypes, strings.ToUpper(fileExt))
	if !found {
		res.Message = fmt.Sprintf("Extensión no permitida, elija un archivo %s", strings.Join(mimeTypes, ", "))
		return res
	}

	res.Success = true
	return res
}

func ClearTempFolder(removeNow bool) error {
    currentTime := time.Now()
    if currentTime.Hour() != 8 || !removeNow {
        return nil
    }

    err := os.RemoveAll("static/temp/")
    if err != nil {
        return err
    }

    err = os.Mkdir("static/temp/",0777)
    if err != nil {
        return err
    }

    return nil
}