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
	// Prepare the files
	htmlPath := "./index.html"
	logoPath := "./logo.png"

	// Open index.html
	htmlFile, err := os.Open(htmlPath)
	if err != nil {
		fmt.Println("Error opening HTML file:", err)
		return
	}
	defer htmlFile.Close()

	// Open logo.png
	logoFile, err := os.Open(logoPath)
	if err != nil {
		fmt.Println("Error opening logo file:", err)
		return
	}
	defer logoFile.Close()

	// Create a buffer and multipart writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Add index.html to the "files" field
	htmlFormFile, err := writer.CreateFormFile("files", filepath.Base(htmlPath))
	if err != nil {
		fmt.Println("Error creating form file for HTML:", err)
		return
	}
	if _, err = io.Copy(htmlFormFile, htmlFile); err != nil {
		fmt.Println("Error copying HTML file:", err)
		return
	}

	// Add logo.png to the same "files" field
	logoFormFile, err := writer.CreateFormFile("files", filepath.Base(logoPath))
	if err != nil {
		fmt.Println("Error creating form file for logo:", err)
		return
	}
	if _, err = io.Copy(logoFormFile, logoFile); err != nil {
		fmt.Println("Error copying logo file:", err)
		return
	}

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
		return
	}

	// Save response body to generated.pdf
	outputFile, err := os.Create("generated.pdf")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	if _, err = io.Copy(outputFile, resp.Body); err != nil {
		fmt.Println("Error saving PDF:", err)
		return
	}

	fmt.Println("PDF saved to generated.pdf")
}
