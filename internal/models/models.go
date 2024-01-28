package models

type Function struct {
	UUID           string `json:"uuid"`
	Name           string `json:"name"`
	CodeLocation   string `json:"codeLocation"`
	RootFsLocation string `json:"rootFsLocation"`
	Status         string `json:"status"`
	Handler        string `json:"handler"`
	Runtime        string `json:"runtime"`
}

type ListFunctionsResponse struct {
	Functions []Function `json:"functions"`
}

type ApiErrorResponse struct {
	Message string `json:"message"`
}

type UploadResponse struct {
	Location string `json:"location"`
}

type UploadErrorResponse struct {
	Message string `json:"message"`
}

type InvokeResponse struct {
	InvocationResponse map[string]interface{} `json:"invocationResponse"`
	StatusCode         int
}
