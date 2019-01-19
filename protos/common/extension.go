package common

import "errors"

// CheckMetadata check the menu metadata
func (menu *Menu) CheckMetadata() error {
	if menu.GetMetadata() == nil {
		return errors.New("Menu' Metadata is nil")
	}

	if menu.Metadata.Id == 0 {
		return errors.New("Menu' Id is nil")
	}

	return nil
}
