package stock

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	Create(CreateRequest) (*Stock, error)
	GetByID(string) (*Stock, error)
	GetByProductID(string) (*Stock, error)
	AddQuantity(string, AddRequest) (*Stock, error)
	SubstractQuantity(string, SubstractRequest) (*Stock, error)
}

type service struct {
	repository      Repository
	quantityService quantity.Service
}

func NewService() (Service, error) {
	repository, err := newRepository()
	if err != nil {
		return nil, err
	}

	quantityService := quantity.NewService()

	return &service{
		repository,
		quantityService,
	}, nil
}

type CreateRequest struct {
	ProductID   primitive.ObjectID `json:"productId" binding:"required"`
	ProductName string             `json:"productName" binding:"required"`
	Quantity    quantity.Quantity  `json:"quantity" binding:"required"`
}

func (s *service) Create(req CreateRequest) (*Stock, error) {
	stock := NewStock(req.ProductID, req.ProductName, req.Quantity)

	stock, err := s.repository.Insert(stock)
	if err != nil {
		return nil, err
	}

	return stock, nil
}

func (s *service) GetByID(id string) (*Stock, error) {
	return s.repository.FindByID(id)
}

func (s *service) GetByProductID(productID string) (*Stock, error) {
	return s.repository.FindByProductID(productID)
}

type AddRequest struct {
	quantity.Quantity
}

func (s *service) AddQuantity(stockId string, req AddRequest) (*Stock, error) {
	stock, err := s.repository.FindByID(stockId)
	if err != nil {
		return nil, err
	}

	newQuantity, err := s.quantityService.Add(&stock.Quantity, &req.Quantity)
	if err != nil {
		return nil, err
	}

	stock.Quantity = *newQuantity

	// TODO: should update
	s.repository.Insert(stock)

	return stock, nil
}

type SubstractRequest struct {
	quantity.Quantity
}

func (s *service) SubstractQuantity(stockId string, req SubstractRequest) (*Stock, error) {
	stock, err := s.repository.FindByID(stockId)
	if err != nil {
		return nil, err
	}

	newQuantity, err := s.quantityService.Substract(&stock.Quantity, &req.Quantity)
	if err != nil {
		return nil, err
	}

	stock.Quantity = *newQuantity

	// TODO: should update
	s.repository.Insert(stock)

	return stock, nil
}
