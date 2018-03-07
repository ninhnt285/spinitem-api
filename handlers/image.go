package handlers

import (
	"encoding/json"
	"net/http"

	imageModel "../models/image"
	ufModel "../models/uploadfile"
)

// UploadImage save image file to disk
var UploadImage http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	http.MaxBytesReader(w, r.Body, 2<<24)
	newFile := ufModel.UploadFile{}
	err := newFile.Add(r, "images/")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, newFile)
}

// AddImage save new image to DB
var AddImage http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Parse image from JSON
	var image imageModel.Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	// Add new Image
	if err := image.Add(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	image.PrepareResult()
	respondWithJSON(w, http.StatusOK, image)
}

// GetAllImages get all images
var GetAllImages http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Get itemId
	r.ParseForm()
	itemID := r.FormValue("item_id")
	if itemID == "" {
		respondWithError(w, http.StatusBadRequest, "Can not found item_id")
		return
	}
	// Get images
	images, err := imageModel.GetAllImages(itemID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	for index := range images {
		images[index].PrepareResult()
	}
	respondWithJSON(w, http.StatusOK, map[string][]imageModel.Image{"images": images})
}
