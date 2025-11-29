package strategy

import "fmt"

// Factory creates a new strategy from param map

type Factory func(params map[string]string) (Strategy, error)

type Registry struct {
	factories map[string]Factory
}

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]Factory),
	}
}

func (r *Registry) Register(id string, factory Factory) {
	if _, exists := r.factories[id]; exists {
		panic(fmt.Sprintf("Strategy %s already registered", id))
	}
	r.factories[id] = factory
}

func (r *Registry) New(id string, params map[string]string) (Strategy, error) {
	factory, ok := r.factories[id]
	if !ok {
		return nil, fmt.Errorf("strategy %s not found", id)
	}
	return factory(params)
}
