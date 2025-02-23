package services

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type AzureManager struct {
	client        *azblob.Client
	containerName string
}

func NewAzureManager() (AzureManager, error) {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	accountKey := os.Getenv("AZURE_STORAGE_KEY")
	blobEndpoint := os.Getenv("AZURE_STORAGE_BLOB_ENDPOINT")
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")

	azureManager := AzureManager{
		containerName: containerName,
	}

	if accountName == "" || accountKey == "" || blobEndpoint == "" {
		return azureManager, fmt.Errorf("AZURE_STORAGE_ACCOUNT, AZURE_STORAGE_KEY, and AZURE_STORAGE_BLOB_ENDPOINT must be set")
	}

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return azureManager, fmt.Errorf("failed to create credential: %w", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(blobEndpoint, credential, nil)
	if err != nil {
		return azureManager, fmt.Errorf("failed to create client: %w", err)
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

	_, err = a.client.UploadFile(ctx, a.containerName, blobName, file, nil)
	if err != nil {
		return fmt.Errorf("failed to upload file to blob: %w", err)
	}

	fmt.Println("File uploaded successfully:", blobName)
	return nil
}

func (a AzureManager) ListBlobs() error {
	pager := a.client.NewListBlobsFlatPager(a.containerName, nil)
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

func (a AzureManager) GetBlobFile(ctx context.Context, blobName, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	_, err = a.client.DownloadFile(ctx, a.containerName, blobName, file, nil)
	if err != nil {
		return fmt.Errorf("downloading blob to file: %w", err)
	}

	return nil
}
