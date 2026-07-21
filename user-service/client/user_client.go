package client

import (
	"net/http"

	"user-service/pb/pbconnect"
)

func NewJobClient(targetURL string) pbconnect.JobServiceClient {
	// targetURL = "http://localhost:50052"
	return pbconnect.NewJobServiceClient(
		http.DefaultClient,
		targetURL,
	)
}