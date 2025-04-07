package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// File to upload
	filePath := "./index.html"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a buffer and multipart writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create the form file field
	formFile, err := writer.CreateFormFile("files", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the file content into the form field
	_, err = io.Copy(formFile, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	// imagePath := "./logo.png"
	// imageFile, err := os.Open(imagePath)
	// if err != nil {
	// 	fmt.Println("Error opening imageFile:", err)
	// 	return
	// }
	// defer imageFile.Close()

	// // Create a buffer and multipart writer
	// var requestImageBody bytes.Buffer
	// writer = multipart.NewWriter(&requestImageBody)

	// Close the writer to finalize the form
	writer.Close()

	// Create the HTTP request
	url := "http://localhost:3000/forms/chromium/convert/html"
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code:", resp.Status)
		fmt.Println("Error: received non-200 response code:", resp.StatusCode)
		return
	}

	// Save response body to my.pdf
	outputFile, err := os.Create("generated.pdf")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		fmt.Println("Error saving PDF:", err)
		return
	}

	fmt.Println("PDF saved to generated.pdf")
}
