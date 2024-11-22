package session

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type storeProvider struct {
	store Store
}

func NewStoreProvider(store Store) StoreProvider {
	return &storeProvider{
		store: store,
	}
}

func (sp *storeProvider) GetStore(_ *http.Request) Store {
	return sp.store
}

func (sp *storeProvider) SetStore(_ echo.Context, store Store) {
	sp.store = store
}
