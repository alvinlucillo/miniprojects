package utils

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DefaultAzureBlobKey     = "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
	DefaultAzureAccountName = "devstoreaccount1"
	TestContainerName       = "files"
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
		return nil, "", fmt.Errorf("starting Azurite container: %w", err)
	}

	// Get Azurite's mapped port
	hostIP, err := azuriteContainer.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("getting container host: %w", err)
	}
	mappedPort, err := azuriteContainer.MappedPort(ctx, "10000")
	if err != nil {
		return nil, "", fmt.Errorf("getting mapped port: %w", err)
	}

	blobEndpoint := fmt.Sprintf("http://%s:%s/devstoreaccount1", hostIP, mappedPort.Port())
	credential, err := azblob.NewSharedKeyCredential(DefaultAzureAccountName, DefaultAzureBlobKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create credential: %w", err)
	}

	blobClient, err := azblob.NewClientWithSharedKeyCredential(blobEndpoint, credential, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create client: %w", err)
	}

	containerClient := blobClient.ServiceClient().NewContainerClient(TestContainerName)
	_, err = containerClient.Create(ctx, &container.CreateOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to create container: %w", err)
	}

	return azuriteContainer, blobEndpoint, nil
}
