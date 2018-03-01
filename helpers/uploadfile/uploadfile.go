package uploadfile

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"../config"
)

// UploadFile save file info
type UploadFile struct {
	Name     string `bson:"name" json:"name"`
	Ext      string `bson:"ext" json:"ext"`
	FilePath string `bson:"file_path" json:"file_path"`
	Size     int64  `bson:"size" json:"size"`
}

// SaveFile save upload file to server
func SaveFile(r *http.Request, fileDir string) (*UploadFile, error) {
	// Read uploadfile
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		return nil, err
	}
	// Open new file in server
	conf := config.GetInstance()
	currentTime := time.Now()
	f, err := os.OpenFile(conf.PublicDir+"/images/"+string(currentTime.Year())+"/"+currentTime.Month().String()+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	// Copy file
	io.Copy(f, file)
	var savedFile = UploadFile{
		Name:     handler.Filename,
		Ext:      filepath.Ext(handler.Filename),
		FilePath: "./test"}
	return &savedFile, nil
}
