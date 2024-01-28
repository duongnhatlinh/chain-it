package api

import "errors"

var ErrSaveBlock = errors.New("Error in saving block")
var ErrCreateEvent = errors.New("Error in creating event")
var ErrGetLastBlock = errors.New("Error in getting last block")
var ErrUndefinedEngineMode = errors.New("Error in consensus type")
