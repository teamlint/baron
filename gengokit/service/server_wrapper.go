package service

import (
	"io"

	"github.com/pkg/errors"
	"github.com/teamlint/baron/gengokit"
)

// NewServerWrapper returns a new server wrapper
func NewServerWrapper(prev io.Reader) gengokit.Renderable {
	return &ServerWrapperRender{
		prev: prev,
	}
}

type ServerWrapperRender struct {
	prev io.Reader
}

// Render will return the existing file if it exists, otherwise it will return
// a brand new copy from the template.
func (r *ServerWrapperRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != ServerWrapperPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (r *ServerWrapperRender) IsFirst() bool {
	return r.prev == nil
}

// IsModified 代码是否更改
func (r *ServerWrapperRender) IsModified() bool {
	return false
}
