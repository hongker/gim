package dto

type Validatable interface {
	Validate() error
}

func Validate(obj any) error {
	if validator, ok := obj.(Validatable); ok {
		return validator.Validate()
	}
	return nil
}
