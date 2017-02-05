package storage

import (
	"github.com/crazyprograms/sud/core"
)

type Storage struct {
	client core.IClient
}

func StartStorage(client core.IClient) *Storage {
	return &Storage{client: client}
}
