package apifunc

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseOK struct {
	Message string `json:"message"`
}

type BookRequest struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type GetAllParams struct {
	Limit  int32  `json:"limit" binding:"required" default:"10"`
	Page   int32  `json:"page" binding:"required" default:"1"`
	Search string `json:"search"`
}


type GetAllCategoriesResult struct {
	Categories []*BookRequest
	Count      int32
}