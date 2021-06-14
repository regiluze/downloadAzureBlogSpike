package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
	"os"
)

const (
	blobURL = "https://froga.blob.core.windows.net/conf/"
)

func main() {
	fmt.Println(" >>>> kaixo ")
	ctx := context.Background()
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT"), os.Getenv("AZURE_STORAGE_ACCESS_KEY")
	if len(accountName) == 0 || len(accountKey) == 0 {
		fmt.Println("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(blobURL)

	containerURL := azblob.NewContainerURL(*URL, p)
	azureBlobUrl := containerURL.NewBlobURL("cloud-connector.yml")
	fmt.Println("?>??>>>>> ", azureBlobUrl)

	downloadResponse, err := azureBlobUrl.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		fmt.Println(">>>> error when reading to buffer ", err)
	}

	fmt.Printf("Downloaded the blob: " + downloadedData.String())

}
