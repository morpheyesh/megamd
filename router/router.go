package router

import (
	"errors"
	"fmt"
	"strings"
)

type routerFactory func(string) (Router, error)

var (
	ErrInvalidCName    = errors.New("CName invalid. eg: megam.megambox.com, www.google.com, typo.com")
	ErrCNameExists     = errors.New("CName already exists")
	ErrCNameNotFound   = errors.New("CName not found")
	ErrCNameNotAllowed = errors.New("CName as router subdomain not allowed")
)

var routers = make(map[string]routerFactory)

// Register registers a new router.
func Register(name string, r routerFactory) {
	routers[name] = r
}

// Get gets the named router from the registry.
func Get(name string) (Router, error) {
	factory, ok := routers[name]
	if !ok {
		return nil, fmt.Errorf("unknown router: %q.", name)
	}
	r, err := factory(name)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Router is the basic interface of this package. It provides methods for
// managing backends and routes. Each backend can have multiple routes.
type Router interface {
	SetCName(cname, name string) error
	UnsetCName(cname, name string) error
	Addr(name string) (string, error)
}

type MessageRouter interface {
	StartupMessage() (string, error)
}

func ChopDomain(cname string) (string, error) {
	sdoms := splitDomainName(cname)
	if sdoms != nil && len(sdoms) >= 2 {
		return strings.Join(sdoms[:len(sdoms)-2], "."), nil
	}
	return "", ErrInvalidCName
}
