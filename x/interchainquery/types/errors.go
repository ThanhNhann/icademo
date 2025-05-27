package types

import "errors"

var (
	ErrAlreadyFulfilled  = errors.New("query already fulfilled")
	ErrInvalidICQRequest = errors.New("invalid ICQ request")
	ErrInvalidICQProof  = errors.New("invalid ICQ proof")
	ErrFailedToRetryQuery = errors.New("failed to retry query")
	ErrICQCallbackNotFound = errors.New("ICQ callback not found")
)
