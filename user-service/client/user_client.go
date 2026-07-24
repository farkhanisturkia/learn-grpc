package client

import (
	"net/http"

	"user-service/pb/pbconnect"
)

func NewJobClient(targetURL string) pbconnect.JobServiceClient {

	return pbconnect.NewJobServiceClient(
		http.DefaultClient,
		targetURL,
	)
}

func NewEducationClient(targetURL string) pbconnect.EducationServiceClient {

	return pbconnect.NewEducationServiceClient(
		http.DefaultClient,
		targetURL,
	)
}