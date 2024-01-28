package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/rs/xid"
	"github.com/runvelocity/cli/internal/models"
	"github.com/runvelocity/cli/utils"
)

type ApiClient struct {
	BaseUrl string
}

// TODO: Update this to return proper response
// func readBody(obj *interface{}, responseBody io.ReadCloser) (*interface{}, error) {
// 	defer responseBody.Close()
// 	bytes, err := io.ReadAll(responseBody)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(bytes, obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return obj, nil
// }

func (c *ApiClient) ListFunctions() (*[]models.Function, error) {

	getRequest, err := http.NewRequest("GET", c.BaseUrl+"/functions", nil)
	if err != nil {
		return nil, err
	}
	getResponse, err := http.DefaultClient.Do(getRequest)
	if err != nil {
		return nil, err
	}
	if getResponse.StatusCode != http.StatusOK {
		var response models.ApiErrorResponse
		defer getResponse.Body.Close()
		bytes, err := io.ReadAll(getResponse.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(response.Message)
	}
	var response models.ListFunctionsResponse
	defer getResponse.Body.Close()
	bytes, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	return &response.Functions, nil
}

func (c *ApiClient) UploadCode(filePath string) (*models.UploadResponse, error) {
	functionId := xid.New().String()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	values := map[string]io.Reader{
		"code": file,
		"key":  strings.NewReader(functionId),
	}
	var uploadBuffer bytes.Buffer
	w := multipart.NewWriter(&uploadBuffer)
	err = utils.WriteBytes(values, w)
	if err != nil {
		return nil, err
	}
	uploadRequest, err := http.NewRequest("POST", c.BaseUrl+"/upload", &uploadBuffer)
	if err != nil {
		return nil, err
	}
	uploadRequest.Header.Set("Content-Type", w.FormDataContentType())
	uploadResponse, err := http.DefaultClient.Do(uploadRequest)
	if err != nil {
		return nil, err
	}
	if uploadResponse.StatusCode != http.StatusOK {
		var response models.UploadErrorResponse
		defer uploadResponse.Body.Close()
		bytes, err := io.ReadAll(uploadResponse.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(response.Message)
	}
	var response models.UploadResponse
	defer uploadResponse.Body.Close()
	bytes, err := io.ReadAll(uploadResponse.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ApiClient) CreateFunction(functionName, filePath, handler, runtime string) (*models.Function, error) {
	uploadResponse, err := c.UploadCode(filePath)
	if err != nil {
		return nil, err
	}

	functionId := xid.New().String()

	// Create the function
	createFunctionArgs := map[string]io.Reader{
		"name":         strings.NewReader(functionName),
		"codeLocation": strings.NewReader(uploadResponse.Location),
		"handler":      strings.NewReader(handler),
		"uuid":         strings.NewReader(functionId),
		"runtime":      strings.NewReader(runtime),
	}

	var createBuffer bytes.Buffer
	w := multipart.NewWriter(&createBuffer)
	err = utils.WriteBytes(createFunctionArgs, w)
	if err != nil {
		return nil, err
	}

	createRequest, err := http.NewRequest("POST", c.BaseUrl+"/functions", &createBuffer)
	if err != nil {
		return nil, err
	}
	createRequest.Header.Set("Content-Type", w.FormDataContentType())
	createResponse, err := http.DefaultClient.Do(createRequest)
	if err != nil {
		return nil, err
	}
	// Check the response
	if createResponse.StatusCode != http.StatusCreated {
		var response models.ApiErrorResponse
		defer createResponse.Body.Close()
		bytes, err := io.ReadAll(createResponse.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bytes, &response)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(response.Message)
	}
	var response models.Function
	defer createResponse.Body.Close()
	bytes, err := io.ReadAll(createResponse.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil

}

func (c *ApiClient) DeleteFunction(functionName string) (*models.Function, error) {
	functions, err := c.ListFunctions()
	if err != nil {
		return nil, err
	}

	for _, v := range *functions {
		if v.Name == functionName {
			deleteRequest, err := http.NewRequest("DELETE", c.BaseUrl+"/functions/"+v.UUID, nil)
			if err != nil {
				return nil, err
			}
			deleteResponse, err := http.DefaultClient.Do(deleteRequest)
			if err != nil {
				return nil, err
			}
			if deleteResponse.StatusCode != http.StatusOK {
				var response models.ApiErrorResponse
				defer deleteResponse.Body.Close()
				bytes, err := io.ReadAll(deleteResponse.Body)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(bytes, &response)
				if err != nil {
					return nil, err
				}
				return nil, errors.New(response.Message)
			}
			var response models.Function
			defer deleteResponse.Body.Close()
			bytes, err := io.ReadAll(deleteResponse.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bytes, &response)
			if err != nil {
				return nil, err
			}
			return &response, nil
		}
	}

	return nil, errors.New("Function not found")

}

func (c *ApiClient) InvokeFunction(functionName string, payload map[string]interface{}) (*models.InvokeResponse, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	functions, err := c.ListFunctions()
	for _, v := range *functions {
		if v.Name == functionName {
			invokeRequest, err := http.NewRequest("POST", c.BaseUrl+"/functions/invoke/"+v.UUID, bytes.NewBuffer(jsonPayload))
			if err != nil {
				return nil, err
			}
			invokeRequest.Header.Set("Content-type", "application/json")
			invokeResponse, err := http.DefaultClient.Do(invokeRequest)
			// Check the response
			if invokeResponse.StatusCode != http.StatusOK {
				var response models.ApiErrorResponse
				defer invokeResponse.Body.Close()
				bytes, err := io.ReadAll(invokeResponse.Body)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(bytes, &response)
				if err != nil {
					return nil, err
				}
				return nil, errors.New(response.Message)
			}
			var response models.InvokeResponse
			defer invokeResponse.Body.Close()
			bytes, err := io.ReadAll(invokeResponse.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bytes, &response)
			if err != nil {
				return nil, err
			}
			response.StatusCode = invokeResponse.StatusCode
			return &response, nil
		}
	}
	return nil, errors.New("Function not found")
}
