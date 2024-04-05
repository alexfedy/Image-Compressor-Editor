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

func mainIndexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
  <head>
    <title>Image Compressor / Editor</title>
    <style>
      :root {
        --blue: #5680e9;
        --lightestblue: #84ceeb;
        --lightblue: #5ab9ea;
        --pale: #c1c8e4;
        --purple: #8860d0;
        --secondary: #6b3ac3;
      }
      * {
        box-sizing: border-box;
      }
      body {
        max-width: 1200px;
        margin: 20px;
        font-family: system-ui;
        display: flex;
        justify-content: center;
        padding: 0;
      }
      form {
        display: flex;
        flex-direction: column;
        align-items: center;
      }
      h1 {
        margin: 10px auto;
      }
      p {
        font-size: small;
        color: rgb(74, 65, 65);
      }
      .formContent {
        display: flex;
        flex-direction: column;
      }
      .formContent > *,
      #advancedOptionsContainer > * {
        margin-bottom: 1rem;
      }
      .formContainer {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: 1rem;
        margin: 10px 1rem;
      }
      .grid {
        display: grid;
        gap: 0.5rem;
      }
      label {
        font-size: small;
        display: inline-block;
        margin-right: 10px;
      }
      input {
        font-size: medium;
        padding: 10px;
        border: 0;
        border-bottom: 2px solid var(--purple);
        background-color: transparent;
      }
      input:focus {
        outline: var(--purple);
      }
      .button {
        width: 100%;
        appearance: none;
        -webkit-appearance: none;
        padding: 10px;
        border: none;
        background-color: var(--purple);
        color: #fff;
        font-weight: 600;
        transition: all 0.3s ease-in;
        box-shadow: 2px 2px 5px var(--purple);
      }
      .advanced {
        width: 100%;
        appearance: none;
        -webkit-appearance: none;
        padding: 5px 10px;
        border: none;
        background-color: transparent;
        color: var(--secondary);
        text-decoration: underline;
      }
      .invisible {
        display: none;
      }
      .visible {
        display: flex;
        flex-direction: column;
      }
      .uploadButton {
        border: 2px solid var(--purple);
        color: var(--purple);
        padding: 5px 10px;
        cursor: pointer;
        transition: all 0.1s ease;
      }
      .uploadButton:hover {
        background-color: var(--purple);
        color: white;
      }
    </style>
  </head>
  <body>
    <form action="/upload" method="post" enctype="multipart/form-data">
      <h1>Image Editor/Compressor</h1>
      <p>Enter information to edit and compress your image</p>
      <div class="grid">
        <button type="button" class="advanced" id="showAdvanced">
          Show advanced options
        </button>
      </div>
      <br />
      <div class="formContainer">
        <div class="formContent">
          <div class="">
            <input
              id="actual-btn"
              type="file"
              name="image"
              accept="image/*"
              required
            />
          </div>
          <div class="grid">
            <label>Size</label>
            <input
              placeholder="width"
              type="number"
              value="100"
              min="1"
              max="1000"
              name="width"
            />
            <span>x</span>
            <input
              placeholder="height"
              type="number"
              value="100"
              min="1"
              max="1000"
              name="height"
            />
          </div>
          <div class="grid">
            <label for="option1"
              ><input
                type="radio"
                id="option1"
                name="choice"
                value="JPEG"
                checked
              />JPEG</label
            >
            <label for="option2">
              <input
                type="radio"
                id="option2"
                name="choice"
                value="PNG"
              />PNG</label
            >
          </div>
        </div>
        <div class="formContent">
          <div id="advancedOptionsContainer" class="invisible">
            <div class="grid">
              <label>Sharpnesss value (number)</label>
              <input
                type="number"
                name="sharpness"
                value="0"
                min="0"
                max="100"
              />
            </div>
            <div class="grid">
              <label>Blur</label>
              <input type="number" name="blur" value="0" min="0" max="100" />
            </div>
            <div class="grid">
              <label>Gamma Correction (1 is default, 0.0 - 3.0 range)</label>
              <input
                type="number"
                name="gamma"
                value="1"
                min="0"
                max="3"
                step="0.1"
              />
            </div>
            <div class="grid">
              <label>Contrast (-100 to 100)</label>
              <input
                type="number"
                name="contrast"
                value="0"
                min="-100"
                max="100"
                step="1"
              />
            </div>
            <div class="grid">
              <label>Brightness (-100 to 100)</label>
              <input
                type="number"
                name="brightness"
                value="0"
                min="-100"
                max="100"
                step="1"
              />
            </div>
            <div class="grid">
              <label>Saturation (-100 to 100)</label>
              <input
                type="number"
                name="saturation"
                value="0"
                min="-100"
                max="100"
                step="1"
              />
            </div>
          </div>
        </div>
      </div>
      <button type="submit" class="button">Confirm</button>
      <br />
    </form>
    <script type="text/javascript">
      let b = document.getElementById("showAdvanced");
      b.addEventListener("click", () => {
        let optionsContainer = document.getElementById(
          "advancedOptionsContainer"
        );
        optionsContainer.classList.toggle("invisible");
        optionsContainer.classList.toggle("visible");
        if (optionsContainer.classList.contains("invisible")) {
          b.innerText = "Show Advanced Options";
        } else {
          b.innerText = "Hide Advanced Options";
        }
      });
    </script>
  </body>
</html>	
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)

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
	} else {
		println("sharpness is 0 or less than 0. No change to sharpness made.")
	}
	if blur > 0 {
		processedImage = Blur(processedImage, blur)
	} else {
		println("blur is 0 or less than 0. No change to blur made.")
	}
	if gamma != 1 {
		processedImage = GammaCorrection(processedImage, gamma)
	} else {
		println("gamma is 1. No change to gamma made.")
	}
	if contrast != 0 {
		processedImage = Contrast(processedImage, contrast)
	} else {
		println("contrast is 0 or less than 0. No change to contrast made.")
	}
	if brightness != 0 {
		processedImage = Brightness(processedImage, brightness)
	} else {
		println("brightness is 0 or less than 0. No change to brightness made.")
	}
	if saturation != 0 {
		processedImage = Saturation(processedImage, saturation)
	} else {
		println("saturation is 0 or less than 0. No change to saturation made.")
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
	http.HandleFunc("/", mainIndexHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/env", envHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	}
	log.Printf("Server listening on http://%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
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
