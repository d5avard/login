package fs

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func NewClient(ctx context.Context) *firestore.Client {
	projectID := "danysavard-ca-327813"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err.Error())
	}
	return client
}
