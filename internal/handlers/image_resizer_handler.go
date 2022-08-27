package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/mrtkmynsndev/image-resizer/internal/models"
	"github.com/mrtkmynsndev/image-resizer/internal/registry"
	"github.com/mrtkmynsndev/image-resizer/internal/response"
	"github.com/mrtkmynsndev/image-resizer/internal/services"
	"github.com/mrtkmynsndev/image-resizer/internal/utils"
	"github.com/pkg/errors"
)

func Create(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request := models.ImageResizeRequest{}
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return response.InternalServerError(err)
	}

	imageUrls := make([]string, len(request.Images))

	imageService := services.NewImageService(len(request.Images))

	for i, img := range request.Images {
		mimeType, err := imageService.GetMimeTypeFromBase64(img.File)
		if err != nil {
			return response.InternalServerError(err)
		}

		if !utils.IsValidMimeType(mimeType) {
			return response.BadRequest(errors.New(fmt.Sprintf("mimeType: %s is not support", mimeType)))
		}

		img, err := imageService.Decode(img.File)
		if err != nil {
			return response.InternalServerError(err)
		}

		resizeFactory := imageService.GetImageResizeFactory(request.ResizeType)
		resizedIamges := resizeFactory.Resize(img, request.ResizeOption)
		imageAsBytes, err := imageService.ImageAsBytes(resizedIamges, "png")
		if err != nil {
			return response.BadRequest(fmt.Errorf("ImageService.ImageAsBytes: %w", err))
		}

		reader := bytes.NewReader(imageAsBytes)

		fileName := fmt.Sprintf("%s.%s", strings.Replace(uuid.New().String(), "-", "", -1), request.SaveOptions)

		out, err := registry.S3Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("BucketName")),
			Key:    aws.String(fileName),
			Body:   reader,
		})

		if err != nil {
			log.Printf("ERROR [%v]", err)
		}

		imageUrls[i] = out.Location
	}

	return response.Ok(imageUrls)
}
