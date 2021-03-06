package yarf

import ()

// MiddlewareHandler interface provides the methods for request filters
// that needs to run before, or after, every request Resource is executed.
type MiddlewareHandler interface {
	PreDispatch(*Context) error
	PostDispatch(*Context) error
}

// Middleware struct is the default implementation of a Middleware and does nothing.
// Users can either implement both methods or composite this struct into their own.
// Both methods needs to be present to satisfy the MiddlewareHandler interface.
type Middleware struct{}

// PreDispatch includes code to be executed before every Resource request.
func (m *Middleware) PreDispatch(c *Context) error {
	return nil
}

// PostDispatch includes code to be executed after every Resource request.
func (m *Middleware) PostDispatch(c *Context) error {
	return nil
}
