package discollect

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrPluginUnregistered = errors.New("discollect: plugin not registered")
	ErrHandlerNotFound    = errors.New("discollect: handler not found for route")
)

// A Registry stores and indexes all available plugins
type Registry struct {
	plugins       []*Plugin
	pluginsByName map[string]*Plugin

	// handlers is immutable after creation
	handlers map[string]map[*regexp.Regexp]Handler
}

// NewRegistry indexes a list of plugins and precomputes the routing table
func NewRegistry(plugins []*Plugin) (*Registry, error) {
	// store pluginsByName for quick lookup
	pluginsByName := make(map[string]*Plugin)
	for _, p := range plugins {
		pluginsByName[p.Name] = p
	}

	// precompile all regexps
	handlers := make(map[string]map[*regexp.Regexp]Handler)
	for _, p := range plugins {
		handlers[p.Name] = make(map[*regexp.Regexp]Handler)
		for route, handler := range p.Routes {
			re, err := regexp.Compile(route)
			if err != nil {
				return nil, fmt.Errorf("registry: regexp did not compile for plugin %s: route %s: %s", p.Name, route, err)
			}
			handlers[p.Name][re] = handler
		}
	}

	return &Registry{
		plugins:       plugins,
		pluginsByName: pluginsByName,
		handlers:      handlers,
	}, nil
}

// Get returns a a plugin by name
func (r *Registry) Get(name string) (*Plugin, error) {
	p, ok := r.pluginsByName[name]
	if !ok {
		return nil, ErrPluginUnregistered
	}
	return p, nil
}

// HandlerFor is the core "router" used to point Tasks to an individual Handler
func (r *Registry) HandlerFor(pluginName string, rawURL string) (Handler, []string, error) {
	p, ok := r.handlers[pluginName]
	if !ok {
		return nil, nil, ErrPluginUnregistered
	}

	for route, handler := range p {
		if route.MatchString(rawURL) {
			return handler, route.FindStringSubmatch(rawURL), nil
		}
	}

	return nil, nil, ErrHandlerNotFound
}
