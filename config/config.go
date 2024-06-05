package config

import (
	"fmt"

	apollo "github.com/go-kratos/kratos/contrib/config/apollo/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/encoding/yaml"
)

func apolloDecoder(value *config.KeyValue, m map[string]interface{}) error {
	var tmp map[string]interface{}
	err := encoding.GetCodec(yaml.Name).Unmarshal(value.Value, &tmp)
	if err != nil {
		return err
	}

	for k, v := range tmp {
		vm, ok := v.(map[string]interface{})
		if ok {
			for vk, vv := range vm {
				m[vk] = vv
			}
			continue
		}
		m[k] = v
	}
	return nil
}

func LoadConfig(v interface{}, opts ...Option) (config.Config, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.filePath == "" && o.apolloEndpoint == "" {
		return nil, fmt.Errorf("file and apollo is empty")
	}

	var configOpts []config.Option
	if o.filePath != "" {
		configOpts = append(configOpts, config.WithSource(file.NewSource(o.filePath)))
	}

	if o.apolloEndpoint != "" {
		configOpts = append(configOpts, config.WithSource(
			apollo.NewSource(
				apollo.WithAppID(o.apolloAppID),
				apollo.WithCluster(o.apolloCluster),
				apollo.WithEndpoint(o.apolloEndpoint),
				apollo.WithNamespace(o.apolloNamespace),
				apollo.WithSecret(o.apolloSecret),
			),
		),
			config.WithDecoder(apolloDecoder))
	}

	c := config.New(configOpts...)

	if err := c.Load(); err != nil {
		return nil, err
	}

	if o.watchers != nil {
		for key, observer := range o.watchers {
			err := c.Watch(key, observer)
			if err != nil {
				return nil, err
			}
		}
	}

	return c, c.Scan(v)
}

func LoadConfigWithWatcher(v interface{}, opts ...Option) (config.Config, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	if o.filePath == "" && o.apolloEndpoint == "" {
		return nil, fmt.Errorf("file and apollo is empty")
	}

	var configOpts []config.Option
	if o.filePath != "" {
		configOpts = append(configOpts, config.WithSource(file.NewSource(o.filePath)))
	}

	if o.apolloEndpoint != "" {
		configOpts = append(configOpts, config.WithSource(
			apollo.NewSource(
				apollo.WithAppID(o.apolloAppID),
				apollo.WithCluster(o.apolloCluster),
				apollo.WithEndpoint(o.apolloEndpoint),
				apollo.WithNamespace(o.apolloNamespace),
				apollo.WithSecret(o.apolloSecret),
			),
		),
			config.WithDecoder(apolloDecoder))
	}

	c := config.New(configOpts...)

	if err := c.Load(); err != nil {
		return nil, err
	}

	if o.watchers != nil {
		for key, observer := range o.watchers {
			err := c.Watch(key, observer)
			if err != nil {
				return nil, err
			}
		}
	}

	return c, c.Scan(v)
}
