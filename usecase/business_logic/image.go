package businesslogic

import "sonit_server/utils/file"

type imageService struct {
	imageUrl string
}

type iImageService interface {
	IsImageValid() bool
}

func InitializeImageService(imageUrl string) iImageService {
	return &imageService{
		imageUrl: imageUrl,
	}
}

// IsImageValid implements iImageService.
func (i *imageService) IsImageValid() bool {
	return file.IsImageByDataURL(i.imageUrl) && file.IsImageByMagicBytes(i.imageUrl) && file.IsImageByDecoding(i.imageUrl)
}
