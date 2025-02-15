package services

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureManager struct {
	client *azblob.Client
}

const (
	ContainerName = "files"
)

func NewAzureManager() (AzureManager, error) {
	azureManager := AzureManager{}
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_KEY")
	blobEndpoint := os.Getenv("AZURE_STORAGE_BLOB_ENDPOINT") // Required for Azurite

	if accountName == "" || accountKey == "" || blobEndpoint == "" {
		return azureManager, fmt.Errorf("AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_KEY, and AZURE_STORAGE_BLOB_ENDPOINT must be set")
	}

	// Create a connection string similar to the default one Azure provides
	connectionString := fmt.Sprintf(
		"DefaultEndpointsProtocol=http;AccountName=%s;AccountKey=%s;BlobEndpoint=%s",
		accountName, accountKey, blobEndpoint,
	)

	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return azureManager, fmt.Errorf("failed to create blob client: %w", err)
	}

	azureManager.client = client
	return azureManager, nil
}

func (a AzureManager) UploadFile(ctx context.Context, blobName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = a.client.UploadFile(ctx, ContainerName, blobName, file, nil)
	if err != nil {
		return fmt.Errorf("failed to upload file to blob: %w", err)
	}

	fmt.Println("File uploaded successfully:", blobName)
	return nil
}

func (a AzureManager) ListBlobs() error {
	pager := a.client.NewListBlobsFlatPager(ContainerName, nil)
	fmt.Println("Blobs in container:")

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to list blobs: %w", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(" -", *blob.Name)
		}
	}

	return nil
}
