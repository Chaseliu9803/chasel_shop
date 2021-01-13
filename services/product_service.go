package services

import (
	"chasel_shop/datamodels"
	"chasel_shop/repositories"
	"fmt"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProduct()([]*datamodels.Product, error)
	DeleteProductByID(int64)bool
	InsertProduct(product *datamodels.Product)(int64,error)
	UpdateProduct(product *datamodels.Product)error
}

type ProductServie struct {
	productReposity repositories.IProduct
}

//初始化函数
func NewProductService(repository repositories.IProduct)IProductService  {
	return &ProductServie{repository}
}

func (p *ProductServie) GetProductByID(productID int64) (*datamodels.Product, error){
	fmt.Println("productID===== ",productID)
	return p.productReposity.SelectByKey(productID)
}

func (p *ProductServie) GetAllProduct()([]*datamodels.Product, error){
	return p.productReposity.SelectAll()
}

func (p *ProductServie) DeleteProductByID(ProductID int64)bool{
	return p.productReposity.Delete(ProductID)
}

func (p *ProductServie) InsertProduct(product *datamodels.Product)(int64,error){
	return p.productReposity.Insert(product)
}

func (p *ProductServie) UpdateProduct(product *datamodels.Product)error{
	return p.productReposity.Update(product)
}





