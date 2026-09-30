package products

type CreateProductRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Barcode     string `json:"barcode"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,min=0"`
	Stock       int    `json:"stock" binding:"min=0"`
	ImageURL    string `json:"image_url"`
	IsAvailable *bool  `json:"is_available"`
}

type UpdateProductRequest struct {
	CategoryID  int64  `json:"category_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Barcode     string `json:"barcode"`
	Description string `json:"description"`
	Price       int64  `json:"price" binding:"required,min=0"`
	Stock       int    `json:"stock" binding:"min=0"`
	ImageURL    string `json:"image_url"`
	IsAvailable *bool  `json:"is_available"`
}
