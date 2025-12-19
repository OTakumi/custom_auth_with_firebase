package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

// NewClient initializes a new Firebase app and returns a Firestore client.
func NewClient(ctx context.Context) (*firestore.Client, error) {
	// The `FIRESTORE_EMULATOR_HOST` environment variable is automatically used by the
	// library to connect to the emulator.
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return client, nil
}
