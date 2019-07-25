package blob

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/rsdmike/app-functions-azure/pkg/azure"
)

type AzureBlobUpload struct {
	accountInfo   azure.AzureAccountInfo
	containerName string
}

func NewBlobUpload(accountInfo azure.AzureAccountInfo, containerName string) AzureBlobUpload {
	blobUpload := AzureBlobUpload{
		accountInfo:   accountInfo,
		containerName: containerName,
	}
	return blobUpload
}

func (bu AzureBlobUpload) ContainerBlobUpload(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, result interface{}) {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := bu.accountInfo.AccountName
	accountKey := bu.accountInfo.AccountKey

	// Create a ContainerURL object to a container where we'll create a blob and its snapshot.
	// Create a BlockBlobURL object to a blob in the container.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/CopiedBlob.bin", accountName, bu.containerName))
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
	}
	blobURL := azblob.NewBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background() // This example uses a never-expiring context

	src, _ := url.Parse("https://cdn2.auth0.com/docs/media/addons/azure_blob.svg")
	startCopy, err := blobURL.StartCopyFromURL(ctx, *src, nil, azblob.ModifiedAccessConditions{}, azblob.BlobAccessConditions{})
	if err != nil {
		log.Fatal(err)
	}

	copyID := startCopy.CopyID()
	copyStatus := startCopy.CopyStatus()
	for copyStatus == azblob.CopyStatusPending {
		time.Sleep(time.Second * 2)
		getMetadata, err := blobURL.GetProperties(ctx, azblob.BlobAccessConditions{})
		if err != nil {
			log.Fatal(err)
		}
		copyStatus = getMetadata.CopyStatus()
	}
	fmt.Printf("Copy from %s to %s: ID=%s, Status=%s\n", src.String(), blobURL, copyID, copyStatus)
	return true, nil
}
