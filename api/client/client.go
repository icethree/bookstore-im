package client

import (
	"bookstore_im/models/model"
	"fmt"
)

type ImClientOption struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type imClient struct{}

var ImClient ImClientServices = &imClient{}

func NewClient(option *ImClientOption) (ImClientServices, error) {
	err := model.Init(fmt.Sprintf("%s:%s", option.Host, option.Port))
	if err != nil {
		return nil, err
	}

	return ImClient, nil
}

type ImClientServices interface {
	B()
	C()
	D()
}

func (i *imClient) B() {}
func (i *imClient) C() {}
func (i *imClient) D() {}
