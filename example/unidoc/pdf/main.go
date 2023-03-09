package main

import (
	"archive/zip"
	"fmt"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"image/jpeg"
	"os"
)

func main() {
	//if len(os.Args) < 3 {
	//	fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
	//	os.Exit(1)
	//}

	inputPath := "C:\\Users\\17805\\Desktop\\深入理解Nginx模块开发与架构解析.pdf"
	outputPath := "C:\\Users\\17805\\Desktop\\1.zip"

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Extracts images and properties of a PDF specified by inputPath.
// The output images are stored into a zip archive whose path is given by outputPath.
func extractImagesToArchive(inputPath, outputPath string) error {
	pdfReader, f, err := model.NewPdfReaderFromFile(inputPath, nil)
	if err != nil {
		return err
	}
	defer f.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipf.Close()
	zipw := zip.NewWriter(zipf)

	totalImages := 0
	for i := 0; i < numPages; i++ {
		fmt.Printf("-----\nPage %d:\n", i+1)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		pextract, err := extractor.New(page)
		if err != nil {
			return err
		}

		pimages, err := pextract.ExtractPageImages(nil)
		if err != nil {
			return err
		}

		fmt.Printf("%d Images\n", len(pimages.Images))
		for idx, img := range pimages.Images {
			fmt.Printf("Image %d - X: %.2f Y: %.2f, Width: %.2f, Height: %.2f\n",
				totalImages+idx+1, img.X, img.Y, img.Width, img.Height)
			fname := fmt.Sprintf("p%d_%d.jpg", i+1, idx)

			gimg, err := img.Image.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(imgf, gimg, &opt)
			if err != nil {
				return err
			}
		}
		totalImages += len(pimages.Images)
	}
	fmt.Printf("Total: %d images\n", totalImages)

	// Make sure to check the error on Close.
	err = zipw.Close()
	if err != nil {
		return err
	}

	return nil
}
