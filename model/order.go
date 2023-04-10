package model

type Order struct {
	ID          string `json:"id"`
	OrderNumber string `json:"order_number"`
	Price       int    `json:"price"`
	Qty         int    `json:"qty"`
	Total       int    `json:"total"`
	CustomerID  string `json:"customer_id"`
	Status      string `json:"status"`
	Remark      string `json:"remark"`
}
