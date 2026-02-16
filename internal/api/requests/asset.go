package requests

import "time"

type CheckoutAssetPayload struct {
	AssetName           string     `json:"assetName" validate:"required,max=100"`
	AssetID             int64      `json:"assetId" validate:"required"`
	UserID              int64      `json:"userId" validate:"required"`
	CheckoutDate        time.Time  `json:"checkoutDate" validate:"required"`
	ExpectedCheckinDate *time.Time `json:"expectedCheckinDate"`
	Notes               string     `json:"notes"`
}
