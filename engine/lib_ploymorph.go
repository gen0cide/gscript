package engine

import (
	"errors"
)

func (e *Engine) RetrievePolymorphicData() ([]byte, error) {
	return []byte{}, errors.New("not implemented for this platform")
}

func (e *Engine) WritePolymorphicData(polymorphicData []byte) error {
	return errors.New("not implemented for this platform")
}
