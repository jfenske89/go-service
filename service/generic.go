package service

type GenericService struct {
	BaseService
}

func NewGenericService() *GenericService {
	return &GenericService{}
}
