package urlformat

import (
	"fmt"
	"strings"
	"sync"
)

var (
	globalFormat Format = &defaultFormat{}
	lock         sync.Mutex
	httpSchema   = "http://"
	httpsSchema  = "https://"
)

type Format interface {
	FullPath(string, ...Option) string
}

func SetFormat(in Format) {
	lock.Lock()
	defer lock.Unlock()
	globalFormat = in
}

func FullPath(image string, opts ...Option) string {
	if image == "" {
		return ""
	}
	return globalFormat.FullPath(image, opts...)
}

var (
	_ Format = (*defaultFormat)(nil)
	_ Format = (*format)(nil)
)

type defaultFormat struct{}

func (d *defaultFormat) FullPath(s string, opts ...Option) string { return s }

type format struct {
	domain string
}

func NewFormat(domain string) Format {
	if !strings.HasSuffix(domain, "/") {
		domain = domain + "/"
	}

	return &format{
		domain: domain,
	}
}

func (f *format) FullPath(image string, opts ...Option) string {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	if !strings.HasPrefix(image, httpSchema) && !strings.HasPrefix(image, httpsSchema) {
		image = strings.TrimPrefix(image, "/")
		image = fmt.Sprintf("%v%v", f.domain, image)
	}

	if o.template != "" && !strings.HasSuffix(image, o.template) {
		image = fmt.Sprintf("%v%v", image, o.template)
	}
	return image
}
