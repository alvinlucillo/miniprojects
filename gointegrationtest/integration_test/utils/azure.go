package utils

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupAzuriteContainer() (testcontainers.Container, string, error) {
	ctx := context.Background()

	// Start Azurite Testcontainer
	req := testcontainers.ContainerRequest{
		Image:        "mcr.microsoft.com/azure-storage/azurite",
		ExposedPorts: []string{"10000/tcp"},
		WaitingFor:   wait.ForLog("Azurite Blob service is successfully listening"),
	}
	azuriteContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start Azurite container: %w", err)
	}

	// Get Azurite's mapped port
	hostIP, err := azuriteContainer.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get container host: %w", err)
	}
	mappedPort, err := azuriteContainer.MappedPort(ctx, "10000")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Construct the connection string
	blobEndpoint := fmt.Sprintf("http://%s:%s/devstoreaccount1", hostIP, mappedPort.Port())
	return azuriteContainer, blobEndpoint, nil
}
