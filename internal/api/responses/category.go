package responses

import "github.com/knr1997/assets-management-apiserver/internal/store"

type CategoryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewCategoryResponse(u *store.Category) CategoryResponse {
	return CategoryResponse{
		ID:          u.ID,
		Name:        u.Name,
		Description: u.Description,
	}
}

func NewCategorysResponse(Categorys []store.Category) []CategoryResponse {
	responses := make([]CategoryResponse, len(Categorys))

	for i := range Categorys {
		responses[i] = NewCategoryResponse(&Categorys[i])
	}

	return responses
}
