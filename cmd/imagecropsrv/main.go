package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

func unixNano(base int) string {
	un := time.Now().UnixNano()
	return strconv.FormatInt(un, 36)
}

func saveFileHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		imgBase64 := r.FormValue("data")

		// remove "data:image/png;base64,"
		imgBase64cleaned := imgBase64[len("data:image/png;base64,"):len(imgBase64)]

		// decode base64 to buffer bytes
		imgBytes, _ := base64.StdEncoding.DecodeString(imgBase64cleaned)

		// convert []byte to image for saving to file
		r := bytes.NewReader(imgBytes)
		img, _, _ := image.Decode(r)

		x := unixNano(36)

		path := filepath.Join(saveDir, fmt.Sprintf("%v.png", x))
		imgFile, err := os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}

		// save to file on your webserver
		png.Encode(imgFile, img)

	}
}

var (
	saveDir string
)

func cropImageHandle(w http.ResponseWriter, r *http.Request) {

	html := `<!DOCTYPE html>
 	<html lang="en">
 	
 	<head>
 	
 		<title>Golang and Cropper JS demo</title>
 	
 		<!-- cropper js -->
 		<link href="https://unpkg.com/cropperjs@1.5.7/dist/cropper.min.css" rel="stylesheet">
 		<script src="https://unpkg.com/cropperjs@1.5.7/dist/cropper.min.js"></script>
		 
 		<!-- jquery -->
 		<script src="https://unpkg.com/jquery@3.5.1/dist/jquery.min.js" integrity="sha256-9/aliU8dGd2tb6OSsuzixeV4y/faTqgFtohetphbbj0=" crossorigin="anonymous"></script>
 	
 	</head><body>
 	
 		<input type="file" name="img[]" class="file-upload-default" id="cropperImageUpload">
 		<div>
 			   <input type="text" class="form-control file-upload-info" disabled="" placeholder="Upload Image">
 			   <button class="file-upload-browse" type="button">Upload</button>
 		</div>
 	
 	<div>
 	
 		<label>Width (px) :</label>
 		<input type="number" value="300" class="img-w" placeholder="Image width">
 		<button class="btn btn-primary crop mb-2 mb-md-0">Crop</button>
 		<a href="javascript:;" class="save">Save</a>
 	
 	
 	</div>
 	
 	<style>
 	w-100 {
 	  max-width: 100%; /* This rule is very important, please do not ignore this! */
 	}
 	</style>
 	
 	<div>
 		<img src="" class="w-100" id="croppingImage" alt="cropper">
 	</div>
 	
 	<div class="col-md-4 ml-auto">
 		<h6>Cropped Image: </h6>
 		<img class="w-100 cropped-img" src="#" alt="">
 	</div>
 	
 	<script type="text/javascript">
 	$(function() {
 		'use strict';
 	  
 		var croppingImage = document.querySelector('#croppingImage'),
 		img_w = document.querySelector('.img-w'),
 		cropBtn = document.querySelector('.crop'),
 		croppedImg = document.querySelector('.cropped-img'),
 		saveBtn = document.querySelector('.save'),
 		upload = document.querySelector('#cropperImageUpload'),
 	
 		cropper = '';
 	
 		$('.file-upload-browse').on('click', function(e) {
 		  var file = $(this).parent().parent().parent().find('.file-upload-default');
 		  file.trigger('click');
 		});
 	   
 		cropper = new Cropper(croppingImage);
 	  
 		// on change show image with crop options
 		upload.addEventListener('change', function (e) {
 		  $(this).parent().find('.form-control').val($(this).val().replace(/C:\\fakepath\\/i, ''));    
 		  if (e.target.files.length) {
 			cropper.destroy();
 			// start file reader
 			const reader = new FileReader();
 			reader.onload = function (e) {
 			  if(e.target.result){
 				croppingImage.src = e.target.result;
 				cropper = new Cropper(croppingImage);
 			  }
 			};
 			reader.readAsDataURL(e.target.files[0]);
 		 }
 		});
 	  
 		// crop on click
 		cropBtn.addEventListener('click',function(e) {
 		  e.preventDefault();
 		  // get result to data uri
 		  let imgSrc = cropper.getCroppedCanvas({
 			width: img_w.value // input value
 			
 	
 		  }).toDataURL();
 	
 		  croppedImg.src = imgSrc;
 		  saveBtn.setAttribute('href', imgSrc);
 		  
 	
 		});
 	
 		// save on click
 		saveBtn.addEventListener('click',function(e) {
 		  e.preventDefault();
 		  // get result to data uri
 		  let imgSrc = cropper.getCroppedCanvas({
 			//width: img_w.value // input value
 			width: 300 // input value
 		  }).toDataURL();
 	
 		  
 		  //post base64 image data to saveCroppedPhoto function then redirect to the /cropped directory show the actual file
 		  $.post("/saveCroppedPhoto", {data: imgSrc}, function () {window.location.href = "/cropped/"});
 		  
 	
 		});
 	  
 	  
 	  });
 	</script>`

	w.Write([]byte(html))
}

func main() {
	workDir, _ := os.Getwd()
	// croppedDir := filepath.Join(workDir, "/cropped")

	addr := flag.String("addr", "localhost:8080", "TCP address to listen to")
	dir := flag.String("dir", workDir, "Directory to save images")
	flag.Parse()

	saveDir = *dir
	// log.Fatalln(saveDir)
	err := createDirIfNotExist(*dir)

	if err != nil {
		log.Fatalln(err)
		return
	}

	mux := chi.NewRouter()
	mux.HandleFunc("/", cropImageHandle)
	mux.HandleFunc("/saveCroppedPhoto", saveFileHandle)

	// expose our /cropped directory

	// log.Fatalln(croppedDir)
	fileServer(mux, "/cropped", http.Dir(*dir))

	err = http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatalln(err)
	}
}

func createDirIfNotExist(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// log.Fatalln("abc")
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
	}

	if !info.IsDir() {
		return errors.New(dir + " is not a directory")
	}
	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
