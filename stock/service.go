package stock

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/aboglioli/big-brother/quantity"
)

type Service interface {
	Create(CreateRequest) (*Stock, error)
	GetByID(string) (*Stock, error)
	GetByProductID(string) (*Stock, error)
	AddQuantity(string, quantity.Quantity) (*Stock, error)
	SubstractQuantity(string, quantity.Quantity) (*Stock, error)
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

func (s *service) AddQuantity(stockID string, q quantity.Quantity) (*Stock, error) {
	stock, err := s.repository.FindByID(stockID)
	if err != nil {
		return nil, err
	}

	newQuantity, err := s.quantityService.Add(&stock.Quantity, &q)
	if err != nil {
		return nil, err
	}

	stock.Quantity = *newQuantity

	// TODO: should update
	s.repository.Insert(stock)

	return stock, nil
}

func (s *service) SubstractQuantity(stockID string, q quantity.Quantity) (*Stock, error) {
	stock, err := s.repository.FindByID(stockID)
	if err != nil {
		return nil, err
	}

	newQuantity, err := s.quantityService.Substract(&stock.Quantity, &q)
	if err != nil {
		return nil, err
	}

	stock.Quantity = *newQuantity

	// TODO: should update
	s.repository.Insert(stock)

	return stock, nil
}
