package repo

import "errors"

var ErrAddBlock = errors.New("Error in adding block")
var ErrGetBlock = errors.New("Error in getting block")
var ErrEmptyBlock = errors.New("Error when block is empty that should be not")
var ErrNewBlockStorage = errors.New("Error in constructing block storage")
