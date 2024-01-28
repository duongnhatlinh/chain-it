package adapter

import "errors"

var ErrBlockNil = errors.New("Block nil error")
var ErrBlockTypeCasting = errors.New("Error in type casting block")
var ErrCommandTransactions = errors.New("command's transactions nil or have length of zero")
var ErrCommandSeal = errors.New("command's transactions nil")
var ErrTxHasMissingProperties = errors.New("Tx has missing properties")
var ErrBlockIdNil = errors.New("Error command model ID is nil")
var ErrTxResultsLengthOfZero = errors.New("Error length of tx results is zero")
var ErrTxResultsFail = errors.New("Error not all tx results success")
var ErrCreateEvent = errors.New("Error in creating consent event")
var ErrBlockSealNil = errors.New("Block seal nil error")
