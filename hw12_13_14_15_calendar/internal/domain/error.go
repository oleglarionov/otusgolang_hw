package domain

import "github.com/pkg/errors"

var ErrAlreadyExists = errors.New("entity already exists")

var ErrNotFound = errors.New("entity not found")
