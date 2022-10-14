package framework

// Engine represents im framework public access api.
type Engine struct{}

// WithSchema use different schema
func (engine *Engine) WithSchema(schema ...Schema) *Engine {
	return engine
}

// WithCallback use callback
func (engine *Engine) WithCallback(callback *Callback) *Engine { return engine }

// WithCodec use codec to pack/unpack message.
func (engine *Engine) WithCodec(codec Codec) *Engine { return engine }

// WithRouter set router
func (engine *Engine) WithRouter(router *Router) *Engine { return engine }

// WithEvent set event
func (engine *Engine) WithEvent(event *Event) *Engine { return engine }

// Start starts the engine
func (engine *Engine) Run() {}

// Close shuts down the engine.
func (engine *Engine) Close() {}

// New returns a new engine instance
func New() *Engine {
	return &Engine{}
}
