package config

import (
	"sync"
)

// Set defines a composite collection of configuration
type Set struct {
	name      string
	path      string
	root      *Set
	parent    *Set
	children  sync.Map
	settings  sync.Map
	notifiers sync.Map
}
