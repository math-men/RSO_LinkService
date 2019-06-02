package models

import "errors"

var	ErrNotFound = errors.New("Requested item is not found!")
var	ErrMarshalling = errors.New("Request body could not be marshalled")
var	ErrInsert = errors.New("Insert operation failed")
var ErrQuery = errors.New("Error during querying")
