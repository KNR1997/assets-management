package responses

import "github.com/knr1997/assets-management-apiserver/internal/store"

type AssetResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serialNumber"`
	Status       string `json:"status"`
	Category     string `json:"category"`
}

func NewAssetResponse(u *store.Asset) AssetResponse {
	return AssetResponse{
		ID:           u.ID,
		Name:         u.Name,
		SerialNumber: u.SerialNumber,
		Status:       string(u.Status),
		Category:     u.Category.Name,
	}
}

func NewAssetsResponse(assets []store.Asset) []AssetResponse {
	responses := make([]AssetResponse, len(assets))

	for i := range assets {
		responses[i] = NewAssetResponse(&assets[i])
	}

	return responses
}
