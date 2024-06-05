package config

import (
	"os"

	"github.com/go-kratos/kratos/v2/config"
)

type options struct {
	filePath        string
	apolloAppID     string
	apolloCluster   string
	apolloEndpoint  string
	apolloNamespace string
	apolloSecret    string
	watchers        map[string]config.Observer
}

type Option func(*options)

func WithWatchers(watchers map[string]config.Observer) Option {
	return func(o *options) {
		o.watchers = watchers
	}
}

func WithFilePath(in string) Option {
	return func(o *options) {
		o.filePath = in
	}
}

func WithAppID(in string) Option {
	return func(o *options) {
		o.apolloAppID = in
	}
}

func WithCluster(in string) Option {
	return func(o *options) {
		o.apolloCluster = in
	}
}

func WithEndpoint(in string) Option {
	return func(o *options) {
		o.apolloEndpoint = in
	}
}

func WithNamespace(in string) Option {
	return func(o *options) {
		o.apolloNamespace = in
	}
}

func WithSecret(in string) Option {
	return func(o *options) {
		o.apolloSecret = in
	}
}

func defaultOptions() *options {
	var (
		cluster    = "default"
		envCluster = os.Getenv("APOLLO_CLUSTER")
	)

	if envCluster != "" {
		cluster = envCluster
	}

	return &options{
		apolloCluster:   cluster,
		apolloAppID:     os.Getenv("APOLLO_APP_ID"),
		apolloEndpoint:  os.Getenv("APOLLO_ENDPOINT"),
		apolloNamespace: os.Getenv("APOLLO_NAMESPACE"),
		apolloSecret:    os.Getenv("APOLLO_SECRET"),
	}
}
