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

func Load(filename string) (*Interface, error) {
	p, err := plugin.Open(filename)
	if err != nil {
		return nil, err
	}

	symExport, err := p.Lookup("Export")
	if err != nil {
		return nil, err
	}

	exportFunc, ok := symExport.(func() Interface)
	if !ok {
		return nil, errors.New("unexpected type from module symbol")
	}

	pluginInstance := exportFunc()

	return &pluginInstance, nil
}
