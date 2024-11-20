package imageprocessor

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // to decode jpeg images
	_ "image/png"  // to decode png images
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
)

const (
	bucketName = "" // give your bucket name
	awsRegion  = ""
)

// DownloadImage downloads the image from a URL and returns the image as bytes.
func DownloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image body: %w", err)
	}

	return imageBytes, nil
}

// CompressImage compresses the image by resizing it.
func CompressImage(imageBytes []byte) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize the image (compressing it)
	compressedImg := resize.Resize(800, 0, img, resize.Lanczos3)

	// Encode the resized image to a byte buffer (you can choose the format here, e.g., JPEG or PNG)
	var buf bytes.Buffer
	err = imaging.Encode(&buf, compressedImg, imaging.JPEG)
	if err != nil {
		return nil, fmt.Errorf("failed to encode compressed image: %w", err)
	}

	return buf.Bytes(), nil
}

// UploadToS3 uploads the compressed image to an S3 bucket and returns the S3 path.
func UploadToS3(imageBytes []byte, fileName string) (string, error) {
	// Create a session to interact with AWS
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := s3.New(sess)

	// Upload the image to S3
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(imageBytes),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	// Return the URL of the uploaded image
	// The URL format will be: https://{bucket-name}.s3.{region}.amazonaws.com/{file-name}
	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, awsRegion, fileName)
	return imageURL, nil
}

// ProcessImage downloads, compresses, and uploads the image to S3
func ProcessImage(imageURL string) (string, error) {
	// 1. Download the image
	log.Printf("Downloading image from URL: %s", imageURL)
	imageBytes, err := DownloadImage(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	// 2. Compress the image
	log.Printf("Compressing image...")
	compressedImageBytes, err := CompressImage(imageBytes)
	if err != nil {
		return "", fmt.Errorf("failed to compress image: %w", err)
	}

	// 3. Generate a unique file name for the compressed image
	fileName := fmt.Sprintf("%s_compressed.jpg", strings.TrimSuffix(imageURL, ".jpg"))

	// 4. Upload to S3
	log.Printf("Uploading compressed image to S3...")
	uploadedImageURL, err := UploadToS3(compressedImageBytes, fileName)
	if err != nil {
		return "", fmt.Errorf("failed to upload compressed image to S3: %w", err)
	}

	// 5. Return the URL of the uploaded image
	return uploadedImageURL, nil
}
