package utils

type Config struct {
	ManagerUrl  string `json:"managerUrl"`
	ManagerPort string `json:"managerPort"`
}

type UploadResponse struct {
	Location string `json:"location"`
	Message  string `json:"message"`
}
