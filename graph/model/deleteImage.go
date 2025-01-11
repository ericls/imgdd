package model

type DeleteImageInput struct {
	ID string `json:"id"`
}

type DeleteImageResult struct {
	ID *string `json:"id,omitempty"`
}
