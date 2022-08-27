package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

var emptyResonse events.APIGatewayProxyResponse = events.APIGatewayProxyResponse{}

type ErrApiResponse struct {
	Message string `json:"message"`
}

func Ok(data interface{}) (events.APIGatewayProxyResponse, error) {
	dataToJson, err := json.MarshalIndent(&data, "", " ")
	if err != nil {
		return emptyResonse, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(dataToJson),
		StatusCode: http.StatusOK,
	}, nil
}

func InternalServerError(err error) (events.APIGatewayProxyResponse, error) {
	log.Printf("ERROR [%s]", err)

	errResponse := ErrApiResponse{
		Message: err.Error(),
	}

	errorToJson, err := json.MarshalIndent(&errResponse, "", " ")
	if err != nil {
		return emptyResonse, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(errorToJson),
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func BadRequest(err error) (events.APIGatewayProxyResponse, error) {
	errResponse := ErrApiResponse{
		Message: err.Error(),
	}

	errorToJson, err := json.MarshalIndent(&errResponse, "", " ")
	if err != nil {
		return emptyResonse, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(errorToJson),
		StatusCode: http.StatusOK,
	}, nil
}

func Created() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
	}, nil
}
