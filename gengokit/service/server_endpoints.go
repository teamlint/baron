package service

import (
	"io"

	"github.com/pkg/errors"
	"github.com/teamlint/baron/gengokit"
)

// NewServerEndpoints returns a new server endpoints render
func NewServerEndpoints(prev io.Reader) gengokit.Renderable {
	return &ServerEndpointsRender{
		prev: prev,
	}
}

type ServerEndpointsRender struct {
	prev io.Reader
}

// Render will return the existing file if it exists, otherwise it will return
// a brand new copy from the template.
func (r *ServerEndpointsRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != ServerEndpointsPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (r *ServerEndpointsRender) IsFirst() bool {
	return r.prev == nil
}

// IsModified 代码是否更改
func (r *ServerEndpointsRender) IsModified() bool {
	return false
}
