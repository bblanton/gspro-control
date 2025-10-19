package actions

import (
	"encoding/json"
	"os"
	"sync"
)

// Combo represents a key combination like ["ctrl", "m"].
type Combo []string

type Registry struct {
	mu      sync.RWMutex
	actions map[string]Combo
}

var (
	registry     *Registry
	registryOnce sync.Once
)

func GetRegistry() *Registry {
	registryOnce.Do(func() {
		registry = &Registry{actions: make(map[string]Combo)}
		registry.load()
	})
	return registry
}

// load reads from actions.json in working directory; falls back to defaults.
func (r *Registry) load() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// defaults
	r.actions = map[string]Combo{
		"mulligan": {"ctrl", "m"},
	}

	f, err := os.Open("actions.json")
	if err != nil {
		return
	}
	defer f.Close()
	var fromFile map[string]Combo
	if err := json.NewDecoder(f).Decode(&fromFile); err == nil {
		for k, v := range fromFile {
			r.actions[k] = v
		}
	}
}

func (r *Registry) List() map[string]Combo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]Combo, len(r.actions))
	for k, v := range r.actions {
		out[k] = v
	}
	return out
}

func (r *Registry) Lookup(name string) (Combo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.actions[name]
	return v, ok
}
