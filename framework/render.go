package framework

type Render struct{}

func (render *Render) Success(data any)  {}
func (render *Render) Failure(err error) {}
