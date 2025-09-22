package request

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name" validate:"required"`
	Description    string `json:"description"`
}

type UpdateCollectionRequest struct {
	CollectionId   string `json:"collection_id" validate:"required"`
	CollectionName string `json:"collection_name"`
	Description    string `json:"description"`
}
