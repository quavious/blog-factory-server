package db

type BadResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
