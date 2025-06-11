package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gen2brain/go-fitz"
)

type Job struct {
	PDFPath string
}

const (
	outputDir  = "output"
	numWorkers = 7
)

// convertPDFtoJPG converts all pages of a PDF to separate JPG images.
func convertPDFtoJPG(pdfPath, outputDir string) error {
	doc, err := fitz.New(pdfPath)
	if err != nil {
		return fmt.Errorf("failed to open PDF %s: %w", pdfPath, err)
	}
	defer doc.Close()

	base := filepath.Base(pdfPath)
	ext := filepath.Ext(base)
	fileName := strings.TrimSuffix(base, ext)

	for i := 0; i < doc.NumPage(); i++ {
		img, err := doc.Image(i)
		if err != nil {
			return fmt.Errorf("failed to render page %d: %w", i, err)
		}

		outPath := filepath.Join(outputDir, fmt.Sprintf("%s-%d.jpg", fileName, i+1))
		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outPath, err)
		}

		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 95})
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to encode image: %w", err)
		}

		outFile.Close()
	}

	return nil
}

// worker continuously pulls jobs from the channel and processes them
func worker(jobs <-chan Job, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	for job := range jobs {
		err := convertPDFtoJPG(job.PDFPath, outputDir)
		if err != nil {
			fmt.Printf("❌ Worker %d: failed to convert %s: %v\n", workerID, job.PDFPath, err)
		} else {
			fmt.Printf("✅ Worker %d: converted %s\n", workerID, job.PDFPath)
		}
	}
}

func main() {
	pdfFiles := []string{"BSE Document.pdf", "Cipla Q3 FY25 Result Update Jan 29 2025 (1).pdf", "DOC July 22 2024.pdf"}

	// Ensure output directory exists
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	jobs := make(chan Job, len(pdfFiles))
	var wg sync.WaitGroup

	// Spawn worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, i)
	}

	// Send jobs (Producer)
	for _, pdfPath := range pdfFiles {
		jobs <- Job{PDFPath: pdfPath}
	}
	close(jobs) // No more jobs

	wg.Wait() // Wait for all workers to finish
}
