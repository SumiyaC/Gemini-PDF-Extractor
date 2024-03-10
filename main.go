//--correct----------

// package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"

// 	//"io/ioutil"
// 	"log"
// 	"net/http"

// 	"cloud.google.com/go/storage"
// 	"cloud.google.com/go/vertexai/genai"
// 	"github.com/gorilla/mux"
// )

// var projectId = "pdfextractor-415909"
// var region = "us-central1"
// var bucketName = "bucket_for_pdf"

// // func extractFromPDF(w http.ResponseWriter, r *http.Request) {
// // 	// Handle file upload (replace with actual error handling)
// // 	file, _, err := r.FormFile("pdf")
// // 	if err != nil {
// // 		http.Error(w, "Error reading uploaded file: "+err.Error(), http.StatusBadRequest)
// // 		return
// // 	}
// // 	defer file.Close()

// func extractFromPDF(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseMultipartForm(32 << 20) // Adjust max memory size as needed
// 	if err != nil {
// 	  http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
// 	  return
// 	}
// 	var pdfBytes []byte  // Declare pdfBytes as a slice of bytes
//     err := io.ReadAll(file, pdfBytes)
// 	if err != nil {
// 		http.Error(w, "Error processing uploaded file: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	// Access uploaded file
// 	file, _, err := r.FormFile("pdf")
// 	if err != nil {
// 	  http.Error(w, "Error reading uploaded file: "+err.Error(), http.StatusBadRequest)
// 	  return
// 	}
// 	defer file.Close()

// 	// Rest of the code remains the same...
// 	var extractedText string
// 	err := io.ReadAll(file, []byte(extractedText)) // Read into existing variable
//   	if err != nil {
//     	http.Error(w, "Error processing uploaded file: "+err.Error(), http.StatusInternalServerError)
//     	return
//   }
//     // Regular expressions for reference and customer number (adjust as needed)
//   refNumRegex := regexp.MustCompile(`(?i)\s*reference number\s*:\s*([^\s]+)`) // Case-insensitive
//   custNumRegex := regexp.MustCompile(`(?i)\s*customer number\s*:\s*([^\s]+)`)

//   // Find reference number using regex
//   refNumMatch := refNumRegex.FindStringSubmatch(extractedText)
//   if refNumMatch != nil {
//     extractedText = refNumMatch[1] // Extract captured group (reference number)
//   } else {
//     extractedText = "Reference number not found"
//   }

//   // Find customer number using regex
//   custNumMatch := custNumRegex.FindStringSubmatch(extractedText)
//   if custNumMatch != nil {
//     extractedText = custNumMatch[1] // Extract captured group (customer number)
//   } else {
//     extractedText = "Customer number not found"
//   }

// 	response := struct {
// 		ReferenceNumber string `json:"referenceNumber"`
// 		CustomerNumber  string `json:"customerNumber"`
// 	}{
// 		ReferenceNumber: extractedText, // Replace with actual extracted reference number
// 		CustomerNumber:  "",            // Replace with actual extracted customer number
// 	}

// 	json.NewEncoder(w).Encode(response)
//   }

//     // var pdfBytes []byte  // Declare pdfBytes as a slice of bytes
//     // err := io.ReadAll(file, pdfBytes)
// 	// if err != nil {
// 	// 	http.Error(w, "Error processing uploaded file: "+err.Error(), http.StatusInternalServerError)
// 	// 	return
// 	// }

// 	// Upload PDF to Cloud Storage (optional, but recommended for scalability)
// 	ctx := context.Background()
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		// Handle Cloud Storage connection error
// 		http.Error(w, "Error connecting to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	writer := client.Bucket(bucketName).Object("uploads/" + r.FormValue("id") + ".pdf").Writer(ctx)
// 	_, err = writer.Write(pdfBytes)
// 	if err != nil {
// 		// Handle Cloud Storage upload error
// 		http.Error(w, "Error uploading PDF to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Call Vertex AI GenAI for text extraction
// 	genaiClient, err := genai.NewClient(ctx, projectId, region)
// 	if err != nil {
// 		// Handle GenAI client creation error
// 		http.Error(w, "Error creating GenAI client: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Replace "your-model-id" with your actual model ID
// 	img := genai.FileData{
// 		MIMEType: "application/pdf", // Adjust for PDF content type
// 		FileURI:  fmt.Sprintf("gs://%s/uploads/%s.pdf", bucketName, r.FormValue("id")),
// 	}
// 	prompt := genai.Text("Extract text from this PDF.")
// 	resp, err := genaiClient.GenerativeModel("projects/"+projectId+"/locations/"+region+"/models/your-model-id").GenerateContent(ctx, img, prompt)
// 	if err != nil {
// 		// Handle GenAI request error
// 		http.Error(w, "Error processing PDF with GenAI: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Parse extracted text (replace with logic to extract reference and customer numbers)

// //}

// func main() {
// 	router := mux.NewRouter()
// 	router.HandleFunc("/extract/{id}", extractFromPDF).Methods(http.MethodPost)
// 	log.Fatal(http.ListenAndServe(":8080", router))
// }
//-------------------

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
	err := r.ParseMultipartForm(32 << 20) // Adjust max memory size as needed
	if err != nil {
		http.Error(w, "Error parsing multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Access uploaded file (optional, for validation purposes)
	file, _, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error reading uploaded file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	var pdfBytes []byte
	pdfBytes, err = io.ReadAll(file) // Assuming 'file' is the uploaded file
	if err != nil {
		http.Error(w, "Error processing uploaded file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Text extraction logic
	ctx := context.Background()

	// Upload PDF to Cloud Storage (optional, but recommended for scalability)
	client, err := storage.NewClient(ctx)
	if err != nil {
		// Handle Cloud Storage connection error
		http.Error(w, "Error connecting to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writer := client.Bucket(bucketName).Object("uploads/" + r.FormValue("id") + ".pdf").NewWriter(ctx)
	_, err = writer.Write(pdfBytes)
	if err != nil {
		// Handle Cloud Storage upload error
		http.Error(w, "Error uploading PDF to Cloud Storage: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Call Vertex AI GenAI for text extraction
	genaiClient, err := genai.NewClient(ctx, projectId, region)
	if err != nil {
		// Handle GenAI client creation error
		http.Error(w, "Error creating GenAI client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Replace "your-model-id" with your actual model ID
	img := genai.FileData{
		MIMEType: "application/pdf", // Adjust for PDF content type
		FileURI:  fmt.Sprintf("gs://%s/uploads/%s.pdf", bucketName, r.FormValue("id")),
	}
	prompt := genai.Text("Extract text from this PDF.")
	resp, err := genaiClient.GenerativeModel("projects/"+projectId+"/locations/"+region+"/models/gemini-1.0-pro:streamGenerateContent").GenerateContent(ctx, img, prompt)
	if err != nil {
		// Handle GenAI request error
		http.Error(w, "Error processing PDF with GenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//-----------
	// Parse extracted text (replace with logic to extract reference and customer numbers based on GenAI response)
	var extractedText string
	// Assuming the response from GenAI is text, modify this based on your model's output format
	//extractedText = resp.GetText()
	log.Printf("GenAI Response: %+v\n", resp)
	extractedText = string(resp.Content())

	// Regular expressions for reference and customer number (adjust as needed)
	refNumRegex := regexp.MustCompile(`(?i)\s*reference number\s*:\s*([^\s]+)`) // Case-insensitive
	custNumRegex := regexp.MustCompile(`(?i)\s*customer number\s*:\s*([^\s]+)`)

	// Find reference number using regex
	refNumMatch := refNumRegex.FindStringSubmatch(extractedText)
	if refNumMatch != nil {
		extractedText = refNumMatch[1] // Extract captured group (reference number)
	} else {
		extractedText = "Reference number not found"
	}

	// Find customer number using regex
	custNumMatch := custNumRegex.FindStringSubmatch(extractedText)
	if custNumMatch != nil {
		extractedText = custNumMatch[1] // Extract captured group (customer number)
	} else {
		extractedText = "Customer number not found"
	}

	// Prepare response
	response := struct {
		ReferenceNumber string `json:"referenceNumber"`
		CustomerNumber  string `json:"customerNumber"`
	}{
		ReferenceNumber: extractedText,
		CustomerNumber:  extractedText, // Use the same variable for now, replace with actual logic for customer number extraction
	}

	// Send response
	json.NewEncoder(w).Encode(response)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/extract", extractFromPDF).Methods("POST")

	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
