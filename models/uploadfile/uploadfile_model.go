package uploadfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/xid"

	"../../helpers/config"
)

// ResizeOption save option of resize process
type ResizeOption struct {
	Width  int
	Height int
	Crop   bool
}

var (
	// ImgOptions save info of all image sizes
	ImgOptions = []ResizeOption{
		// Square images
		ResizeOption{Width: 50, Height: 50, Crop: true},
		ResizeOption{Width: 100, Height: 100, Crop: true},
		ResizeOption{Width: 150, Height: 150, Crop: true},
		ResizeOption{Width: 500, Height: 500, Crop: true},
		// Resize images
		ResizeOption{Width: 500, Height: 500},
		ResizeOption{Width: 750, Height: 750},
		ResizeOption{Width: 1000, Height: 1000}}
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
	f, err := os.OpenFile(cf.PublicDir+destinationDir+fileName+"_origin"+ext, os.O_WRONLY|os.O_CREATE, 0666)
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
	rotateImage(cf.PublicDir+destinationDir+fileName+"_origin"+ext, cf.PublicDir+destinationDir+fileName+ext)
	fmt.Println("2")
	resizeImages(cf.PublicDir+destinationDir, fileName, ext)

	return nil
}

func getOrientation(filePath string) byte {
	var args = []string{filePath, "-f", "exif-ifd0-Orientation"}

	path, err := exec.LookPath("vipsheader")
	if err != nil {
		return 1
	}
	cmd, err := exec.Command(path, args...).Output()
	if err != nil {
		return 1
	}
	return cmd[0] - 48
}

func rotateImage(originFile string, newFile string) error {
	var args = []string{"autorot", originFile, newFile}

	path, err := exec.LookPath("vips")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args...)
	err = cmd.Run()
	fmt.Println("1")
	if err != nil {
		return err
	}
	return err
}

func resizeImages(dir string, filename string, ext string) {
	fmt.Println("3")
	for _, option := range ImgOptions {
		go resizeImage(dir, filename, ext, option)
	}
}

func resizeImage(dir string, filename string, ext string, option ResizeOption) error {
	// Generate new filename
	suffix := "_" + strconv.Itoa(option.Width) + "x" + strconv.Itoa(option.Height)
	if option.Crop {
		suffix += "_square"
	}
	// Generate command
	var args = []string{
		dir + filename + ext,
		"--size", strconv.Itoa(option.Width) + "x" + strconv.Itoa(option.Height),
		"--output", filename + suffix + ext,
	}
	if option.Crop {
		args = append(args, "--crop")
	}

	path, err := exec.LookPath("vipsthumbnail")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args...)
	err = cmd.Run()
	if err != nil {
		return err
	}
	return err
}
