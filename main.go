package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"
)

func bytesToKb(bytes int) float64 {
	return float64(bytes) / 1024
}

func envHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, os.Getenv("TEST_ENV"))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer file.Close()

	imgData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	width := r.FormValue("width")
	height := r.FormValue("height")
	sharpness := r.FormValue("sharpness")
	blur := r.FormValue("blur")
	gamma := r.FormValue("gamma")
	contrast := r.FormValue("contrast")
	brightness := r.FormValue("brightness")
	saturation := r.FormValue("saturation")

	encodingType := r.FormValue("choice")

	width_int, _ := strconv.Atoi(width)
	height_int, _ := strconv.Atoi(height)
	sharpness_float, _ := strconv.ParseFloat(sharpness, 64)
	blur_float, _ := strconv.ParseFloat(blur, 64)
	gamma_float, _ := strconv.ParseFloat(gamma, 64)
	contrast_float, _ := strconv.ParseFloat(contrast, 64)
	brightness_float, _ := strconv.ParseFloat(brightness, 64)
	saturation_float, _ := strconv.ParseFloat(saturation, 64)

	processedImage, originalImageSize, err := processImage(imgData, width_int, height_int, sharpness_float, blur_float, gamma_float, contrast_float, brightness_float, saturation_float, encodingType)
	if err != nil {
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	fileSize := bytesToKb(len(processedImage))

	base64Image := base64.StdEncoding.EncodeToString(processedImage)

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Image Resize/Compress</title>
		</head>
		<body>
			<a href="/">Back</a>
			<p>Encoding: %s</p>
			<p>Original Size: %.2f KB</p>
			<p>New File Size: %.2f KB</p>
			<img src="data:image/jpeg;base64,%s" alt="Embedded Image">
			<p><a href="data:image/jpeg;base64,%s" download="image.jpg"><button>Download Image</button></a></p>
		</body>
		</html>
	`, encodingType, originalImageSize, fileSize, base64Image, base64Image)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func processImage(imgData []byte, width int, height int, sharpness float64, blur float64, gamma float64, contrast float64, brightness float64, saturation float64, encodingType string) ([]byte, float64, error) {

	OG_IMG_SIZE := bytesToKb(len(imgData))

	srcImage, err := imaging.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, OG_IMG_SIZE, err
	}

	processedImage := Resize(srcImage, width, height)
	if sharpness > 0 {
		processedImage = Sharpen(processedImage, sharpness)
	}
	if blur > 0 {
		processedImage = Blur(processedImage, blur)
	}
	if gamma != 1 {
		processedImage = GammaCorrection(processedImage, gamma)
	}
	if contrast != 0 {
		processedImage = Contrast(processedImage, contrast)
	}
	if brightness != 0 {
		processedImage = Brightness(processedImage, brightness)
	}
	if saturation != 0 {
		processedImage = Saturation(processedImage, saturation)
	}

	var buf bytes.Buffer
	if encodingType == "JPEG" {
		if err := imaging.Encode(&buf, processedImage, imaging.JPEG); err != nil {
			return nil, OG_IMG_SIZE, err
		}
	} else {
		if err := imaging.Encode(&buf, processedImage, imaging.PNG); err != nil {
			return nil, OG_IMG_SIZE, err
		}

	}

	return buf.Bytes(), OG_IMG_SIZE, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/env", envHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	}
	log.Printf("Server listening on http://%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
	// log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func Resize(srcImage image.Image, width int, height int) image.Image {
	return imaging.Resize(srcImage, width, height, imaging.Lanczos)
}
func Blur(srcImage image.Image, value float64) image.Image {
	return imaging.Blur(srcImage, value)
}
func Sharpen(srcImage image.Image, value float64) image.Image {
	return imaging.Sharpen(srcImage, value)
}
func GammaCorrection(srcImage image.Image, value float64) image.Image {
	return imaging.AdjustGamma(srcImage, value)
}
func Contrast(srcImage image.Image, value float64) image.Image {
	return imaging.AdjustContrast(srcImage, value)
}
func Brightness(srcImage image.Image, value float64) image.Image {
	return imaging.AdjustBrightness(srcImage, value)
}
func Saturation(srcImage image.Image, value float64) image.Image {
	return imaging.AdjustSaturation(srcImage, value)
}
