package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mrtkmynsndev/image-resizer/internal/handlers"
)

func main() {
	lambda.Start(handlers.Create)
}
