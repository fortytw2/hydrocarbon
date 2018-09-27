package discollect

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrPluginUnregistered         = errors.New("discollect: plugin not registered")
	ErrHandlerNotFound            = errors.New("discollect: handler not found for route")
	ErrNoValidPluginForEntrypoint = errors.New("discollect: no plugin found for entrypoint")
)

// A Registry stores and indexes all available plugins
type Registry struct {
	plugins       []*Plugin
	pluginsByName map[string]*Plugin

	// used to determine what plugin to map to a route
	entrypoints map[string][]*regexp.Regexp

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
	entrypoints := make(map[string][]*regexp.Regexp)
	for _, p := range plugins {
		handlers[p.Name] = make(map[*regexp.Regexp]Handler)
		for route, handler := range p.Routes {
			re, err := regexp.Compile(route)
			if err != nil {
				return nil, fmt.Errorf("registry: regexp did not compile for plugin %s: route %s: %s", p.Name, route, err)
			}
			handlers[p.Name][re] = handler
		}

		entrypoints[p.Name] = make([]*regexp.Regexp, 0)
		for _, e := range p.Entrypoints {
			re, err := regexp.Compile(e)
			if err != nil {
				return nil, fmt.Errorf("registry: entrypoint regexp did not compile for plugin %s: entrypoint %s: %s", p.Name, e, err)
			}

			entrypoints[p.Name] = append(entrypoints[p.Name], re)
		}
	}

	return &Registry{
		plugins:       plugins,
		entrypoints:   entrypoints,
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

// PluginFor finds the
func (r *Registry) PluginFor(entrypointURL string, blacklistNames []string) (*Plugin, []string, error) {
	for _, p := range r.plugins {

		var next = false
		for _, b := range blacklistNames {
			if p.Name == b {
				next = true
			}
		}
		if next {
			continue
		}

		for _, re := range r.entrypoints[p.Name] {
			if re.MatchString(entrypointURL) {
				return p, re.FindStringSubmatch(entrypointURL), nil
			}
		}
	}

	return nil, nil, ErrNoValidPluginForEntrypoint
}
