package utils

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"os"
)

func NewCloudinary() *cloudinary.Cloudinary {
	cld, _ := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	return cld
}
