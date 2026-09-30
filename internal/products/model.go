package products

type Product struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"category_id"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	ImageURL    string `json:"image_url"`
	IsAvailable bool   `json:"is_available"`
}
