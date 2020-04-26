package fooapi

import (
	"fmt"

	"github.com/rancher/steve/pkg/schemaserver/types"
)

type Store struct {
	types.Store
}

func Wrap(store types.Store) types.Store {
	return &Store{
		store,
	}
}

func (s *Store) ByID(apiOp *types.APIRequest, schema *types.APISchema, id string) (types.APIObject, error) {
	return s.Store.ByID(apiOp, schema, id)
}

func (s *Store) Create(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject) (types.APIObject, error) {
	fmt.Printf("call create", apiOp)
	return s.Store.Create(apiOp, schema, data)
}

func (s *Store) Update(apiOp *types.APIRequest, schema *types.APISchema, data types.APIObject, id string) (types.APIObject, error) {
	return s.Store.Update(apiOp, schema, data, id)
}
