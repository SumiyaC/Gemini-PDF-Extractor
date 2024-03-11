package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	//"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/vertexai/genai"
	"github.com/gorilla/mux"
)

var projectId = "pdfextractor-415909"
var region = "us-central1"
var bucketName = "bucket_for_pdf"

func extractFromPDF(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error reading uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	var pdfBytes []byte
	pdfBytes, err = io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error processing uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		http.Error(w, "Error connecting to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer := client.Bucket(bucketName).Object("uploads/" + r.FormValue("id") + ".pdf").NewWriter(ctx)
	_, err = writer.Write(pdfBytes)
	if err != nil {
		http.Error(w, "Error uploading PDF to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	genaiClient, err := genai.NewClient(ctx, projectId, region)
	if err != nil {
		http.Error(w, "Error creating GenAI client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	img := genai.FileData{
		MIMEType: "application/pdf",
		FileURI:  fmt.Sprintf("gs://%s/uploads/%s.pdf", bucketName, r.FormValue("id")),
	}
	prompt := genai.Text("Extract text from this PDF.")
	resp, err := genaiClient.GenerativeModel("projects/"+projectId+"/locations/"+region+"/models/gemini-1.0-pro:streamGenerateContent").GenerateContent(ctx, img, prompt)
	if err != nil {

		http.Error(w, "Error processing PDF with GenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var extractedText string
	log.Printf("GenAI Response: %+v\n", resp)
	extractedText = string(resp.Content())

	refNumRegex := regexp.MustCompile(`(?i)\s*reference number\s*:\s*([^\s]+)`)
	custNumRegex := regexp.MustCompile(`(?i)\s*customer number\s*:\s*([^\s]+)`)

	refNumMatch := refNumRegex.FindStringSubmatch(extractedText)
	if refNumMatch != nil {
		extractedText = refNumMatch[1]
	} else {
		extractedText = "Reference number not found"
	}

	custNumMatch := custNumRegex.FindStringSubmatch(extractedText)
	if custNumMatch != nil {
		extractedText = custNumMatch[1]
	} else {
		extractedText = "Customer number not found"
	}

	response := struct {
		ReferenceNumber string `json:"referenceNumber"`
		CustomerNumber  string `json:"customerNumber"`
	}{
		ReferenceNumber: extractedText,
		CustomerNumber:  extractedText,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/extract", extractFromPDF).Methods("POST")

	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
