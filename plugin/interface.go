package plugin

import (
	"context"
	"errors"
	"plugin"
)

type Interface interface {
	Init(chan IPCMessage) chan IPCMessage
	Attach(context.Context)
	Unload() int
}

type Export func() *Interface

func Load(filename string) (*Interface, error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	symExport, err := p.Lookup("Export")
	if err != nil {
		return nil, err
	}

	exportFunc, ok := symExport.(Export)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}

	pluginInstance := exportFunc()

	return pluginInstance, nil
}
