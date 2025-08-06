package queue

import (
	"fmt"
	"strings"
)

type Factory func(rawURL string) (IBroker, error)

var queues = make(map[string]Factory)

func Register(scheme string, factory Factory) {
	queues[scheme] = factory
}

func ListAvailable() []string {
	var available []string
	for k := range queues {
		available = append(available, k)
	}
	return available
}

func Get(scheme string) (Factory, error) {
	factory, ok := queues[strings.ToLower(scheme)]
	if !ok {
		return nil, fmt.Errorf("unsupported broker type: %s", scheme)
	}
	return factory, nil
}

func init() {
	Register("inmemory", func(rawURL string) (IBroker, error) {
		return &sMemoryBroker{}, nil
	})
}
