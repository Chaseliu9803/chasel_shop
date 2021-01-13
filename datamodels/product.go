package datamodels

type Product struct {
	ID           int64  `json:"id" sql:"ID" chasel_shop:"id"`
	ProductName  string `json:"ProductName" sql:"productName" chasel_shop:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" chasel_shop:"ProductNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" chasel_shop:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" chasel_shop:"ProductUrl"`
}
