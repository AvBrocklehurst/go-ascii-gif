package asciigif

import (
	"fmt"
)

var (
	//ErrInvalidGif is the error returned when an invalid gif is passed to the library
	ErrInvalidGif = fmt.Errorf("Invalid Gif passed")
)
