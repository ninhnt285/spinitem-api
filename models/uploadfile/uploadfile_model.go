package uploadfile

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/xid"

	"../../helpers/config"
)

// UploadFile save file info
type UploadFile struct {
	Title       string `bson:"title" json:"title"`
	Destination string `bson:"destination" json:"destination"`
	Size        int64  `bson:"size" json:"size"`
}

// Add upload file to server
func (uf *UploadFile) Add(r *http.Request, fileDir string) error {
	// Read uploadfile
	file, handler, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	// Open new file in server
	cf := config.GetInstance()
	currentTime := time.Now()
	// Generate Title (original file name) and filePath
	ext := filepath.Ext(handler.Filename)
	destinationDir := fileDir +
		strconv.Itoa(currentTime.Year()) + "/" +
		strconv.Itoa(int(currentTime.Month())+1) + "/" +
		strconv.Itoa(currentTime.Day()) + "/"
	fileName := xid.New().String() + ext

	uf.Title = handler.Filename
	uf.Destination = destinationDir + fileName
	uf.Size = handler.Size
	// Create directory tree
	os.MkdirAll(cf.PublicDir+destinationDir, 0755)
	f, err := os.OpenFile(cf.PublicDir+uf.Destination, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	// Copy file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	return nil
}
