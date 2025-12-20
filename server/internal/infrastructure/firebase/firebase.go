package firebase

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

// NewClient initializes a new Firebase app and returns a Firestore client.
func NewClient(ctx context.Context) (*firestore.Client, *auth.Client, error) {
	// The `FIRESTORE_EMULATOR_HOST` environment variable is automatically used by the
	// library to connect to the emulator.
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new firebase app: %w", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	return firestoreClient, authClient, nil
}
