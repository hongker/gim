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

func ValidateFunc(obj any) func() error {
	return func() error {
		return Validate(obj)
	}
}
