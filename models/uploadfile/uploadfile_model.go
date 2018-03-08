package uploadfile

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/h2non/bimg.v1"

	"github.com/rs/xid"

	"../../helpers/config"
)

var (
	imgOptions = []bimg.Options{
		// Square images
		bimg.Options{Width: 50, Height: 50, Crop: true},
		bimg.Options{Width: 100, Height: 100, Crop: true},
		bimg.Options{Width: 150, Height: 150, Crop: true},
		bimg.Options{Width: 500, Height: 500, Crop: true},
		// Resize images
		bimg.Options{Width: 500, Height: 500, Enlarge: true},
		bimg.Options{Width: 750, Height: 750, Enlarge: true},
		bimg.Options{Width: 1000, Height: 1000, Enlarge: true}}
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
		strconv.Itoa(int(currentTime.Month())) + "/" +
		strconv.Itoa(currentTime.Day()) + "/"
	fileName := xid.New().String()

	uf.Title = handler.Filename
	uf.Destination = destinationDir + fileName + "%s" + ext
	uf.Size = handler.Size
	// Create directory tree
	os.MkdirAll(cf.PublicDir+destinationDir, 0755)
	f, err := os.OpenFile(cf.PublicDir+destinationDir+fileName+ext, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	// Copy file
	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}
	// Resize images
	resizeImages(cf.PublicDir+destinationDir, fileName, ext)

	return nil
}

func resizeImages(dir string, filename string, ext string) {
	for _, option := range imgOptions {
		resizeImage(dir, filename, ext, option)
	}
}

func resizeImage(dir string, filename string, ext string, option bimg.Options) error {
	// Read Image
	buffer, err := bimg.Read(dir + filename + ext)
	if err != nil {
		return err
	}
	// Generate new filename
	suffix := "_" + strconv.Itoa(option.Width) + "x" + strconv.Itoa(option.Height)
	if option.Crop {
		suffix += "_square"
	}
	var newImageName = dir + filename + suffix + ext
	// Resize Image
	newImage, err := bimg.NewImage(buffer).Process(option)
	bimg.Write(newImageName, newImage)
	return err
}
