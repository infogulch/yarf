package yarf

import (
	"net/http"
)

// Framework version string
const Version = "0.6"

// yarf is the main entry point for the framework and it centralizes most of the functionality.
// All configuration actions are handled by this object.
type yarf struct {
	// UseCache indicates if the route cache should be used.
	UseCache bool

	// Debug enables/disables the debug mode.
	// On debug mode, extra error information is sent to the client.
	Debug bool

	// Silent mode attempts to prevent all messages that aren't part of a resource response to get to the client.
	// Specially useful to hide error messages.
	Silent bool

	// PanicHandler can store a func() that will be defered by each request to be able to recover().
	// If you need to log, send information or do anything about a panic, this is your place.
	PanicHandler func()

	GroupRouter

	// Cached routes storage
	cache *Cache
}

// New creates a new yarf and returns a pointer to it.
// Performs needed initializations
func New() *yarf {
	y := new(yarf)

	// Init cache
	y.UseCache = true
	y.cache = NewCache()
	y.GroupRouter = RouteGroup("")

	// Return object
	return y
}

// ServeHTTP Implements http.Handler interface into yarf.
// Initializes a Context object and handles middleware and route actions.
// If an error is returned by any of the actions, the flow is stopped and a response is sent.
func (y *yarf) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if y.PanicHandler != nil {
		defer y.PanicHandler()
	}

	// Set initial context data.
	// The Context pointer will be affected by the middleware and resources.
	c := NewContext(req, res)

	// Cached routes
	if y.UseCache {
		if cache, ok := y.cache.Get(req.URL.Path); ok {
			// Set context params
			c.Params = cache.params
			c.groupDispatch = cache.route

			// Dispatch and stop
			err := y.Dispatch(c)
			y.errorHandler(err, c)
			return
		}
	}

	// Route match
	if y.Match(req.URL.Path, c) {
		if y.UseCache {
			y.cache.Set(req.URL.Path, routeCache{c.groupDispatch, c.Params})
		}
		err := y.Dispatch(c)
		y.errorHandler(err, c)
		return
	}

	// Return 404
	c.Response.WriteHeader(404)
}

// errorHandler deals with request errors.
func (y *yarf) errorHandler(err error, c *Context) {
	// Return if no error or silent mode
	if err == nil || y.Silent {
		return
	}

	// Check error type
	yerr, ok := err.(YError)
	if !ok {
		// Create default 500 error
		yerr = &CustomError{
			httpCode:  500,
			errorCode: 0,
			errorMsg:  err.Error(),
			errorBody: err.Error(),
		}
	}

	// Write error data to response.
	c.Response.WriteHeader(yerr.Code())

	if y.Debug {
		c.Response.Write([]byte(yerr.Body()))
	}
}

// Start initiates a new http yarf server and start listening.
// It's a shortcut for http.ListenAndServe(address, y)
func (y *yarf) Start(address string) {
	http.ListenAndServe(address, y)
}

// StartTLS initiates a new http yarf server and starts listening to HTTPS requests.
// It is a shortcut for http.ListenAndServeTLS(address, cert, key, yarf)
func (y *yarf) StartTLS(address, cert, key string) {
	http.ListenAndServeTLS(address, cert, key, y)
}
