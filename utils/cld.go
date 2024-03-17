package utils

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cloudinaryInstance *cloudinary.Cloudinary

func InitializeCloudinary(connectionCh chan<- string) {
	cloudName := os.Getenv("CLD_CLOUD_NAME")
	apiKey := os.Getenv("CLD_API_KEY")
	apiSecret := os.Getenv("CLD_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatal("Empty cld configiration")
	}

	connectionCh <- "Connecting to cloudinary..."

	cld, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)

	connectionCh <- "Pinging cloudinary..."
	_, err := cld.Admin.Ping(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	connectionCh <- "Connnected to cloudinary"

	cloudinaryInstance = cld
}

/*
Uploads file to cloudinary and returns upload result or error
*/
func UploadFile(file io.Reader) (*uploader.UploadResult, error) {
	return cloudinaryInstance.Upload.Upload(context.Background(), file, uploader.UploadParams{
		PublicIDPrefix: "gin-basic-api",
		ResourceType:   "auto",
	})
}

func DeleteFile(publicId string, resourceType api.AssetType) (*uploader.DestroyResult, error) {
	return cloudinaryInstance.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID:     publicId,
		Invalidate:   api.Bool(true),
		ResourceType: resourceType.String(),
	})
}
