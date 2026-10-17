package msgerror

import "errors"

var (
	ErrNoRows = errors.New("no row affected")
)
