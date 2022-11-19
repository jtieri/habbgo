/*
Package query provides a mechanism for requesting data from a service. Query can also be used to request an update
to some state stored in a running service's cache.

There are two main scenarios that we are accounting for.
	-We are trying to update something and only need a response to signal task completion.
	-We are trying to query for some data that we expect a response back.

-Key
	-What we are querying for? e.g. roomID, playerID, categoryID, etc.
	-The key for some value that we want to update in the running service's cache.

-Value
	-What type of data we are expecting in response to a query.
	-The value we want to update.
*/
package query

// Query describes a Key quantifier that we are querying a running service for and a Value that is expected in return
// or that we are using to update some value cached in the service.
type Query[K, V any] struct {
	Key   K
	Value V
}

// newQuery instantiates a new Query object.
func newQuery[K, V any](key K, value V) *Query[K, V] {
	return &Query[K, V]{
		Key:   key,
		Value: value,
	}
}

// Request contains a Query that is to be executed on some running service along with a Response channel that it
// expects a response on.
type Request[K, V any] struct {
	Query    *Query[K, V]
	Response chan *Query[K, V]
}

// NewRequest instantiates a new Request object.
func NewRequest[K, V any](key K, value V) *Request[K, V] {
	return &Request[K, V]{
		Query:    newQuery(key, value),
		Response: make(chan *Query[K, V]),
	}
}
