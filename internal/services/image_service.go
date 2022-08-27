package services

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/mrtkmynsndev/image-resizer/internal/models"
	"github.com/mrtkmynsndev/image-resizer/internal/utils"
)

type ImageResizeInferface interface {
	Resize(image image.Image, option models.ImageResizeOption) image.Image
}

type ImageResizePixel struct {
}

func (r *ImageResizePixel) Resize(img image.Image, option models.ImageResizeOption) image.Image {
	convertedImage := imaging.Resize(img, option.Width, option.Height, imaging.Lanczos)

	return convertedImage
}

type ImageResizePercentage struct {
}

func (r *ImageResizePercentage) Resize(img image.Image, option models.ImageResizeOption) image.Image {
	log.Printf("image -> %v \n", img)
	width := (1 - option.Percentage) * img.Bounds().Dx()
	height := (1 - option.Percentage) * img.Bounds().Dy()
	convertedImage := imaging.Resize(img, width, height, imaging.Lanczos)

	return convertedImage
}

type ImageServiceInterface interface {
	GetMimeTypeFromBase64(file string) (string, error)
	Decode(file string) (image.Image, error)
	GetImageResizeFactory(resize string) ImageResizeInferface
	ImageAsBytes(img image.Image, imageType string) ([]byte, error)
}

type ImageService struct {
	Images []image.Image
}

// GetImageResizeFactory implements ImageServiceInterface
func (i *ImageService) GetImageResizeFactory(resize string) ImageResizeInferface {
	switch resize {
	case "pixel":
		return &ImageResizePixel{}
	case "percentage":
		return &ImageResizePercentage{}
	default:
		panic("unsupported resize type")
	}
}

// Decode implements ImageServiceInterface
func (i *ImageService) Decode(file string) (image.Image, error) {
	coI := strings.Index(file, ",")
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(file[coI+1:]))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// GetMimeType implements ImageServiceInterface
func (i *ImageService) GetMimeTypeFromBase64(file string) (string, error) {
	mimeType, err := utils.DetectMimeTypeFromBase64(file)
	if err != nil {
		return "", err
	}

	return mimeType, nil
}

// ImageAsBytes implements ImageServiceInterface
func (i *ImageService) ImageAsBytes(img image.Image, imageType string) ([]byte, error) {
	var err error
	buff := new(bytes.Buffer)
	switch imageType {
	case "png":
		err = png.Encode(buff, img)
	case "jpg":
		err = jpeg.Encode(buff, img, nil)
	case "gif":
		err = gif.Encode(buff, img, nil)
	default:
		return nil, errors.New("ImageAsBytes: unsupported image type")
	}

	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// New Image Service
func NewImageService(capacity int) ImageServiceInterface {
	return &ImageService{
		Images: make([]image.Image, capacity),
	}
}
