package service

import (
	"io"

	"github.com/pkg/errors"
	"github.com/teamlint/baron/gengokit"
)

// NewServerInterrupt returns a new server interrupt render
func NewServerInterrupt(prev io.Reader) gengokit.Renderable {
	return &ServerInterruptRender{
		prev: prev,
	}
}

type ServerInterruptRender struct {
	prev io.Reader
}

// Render will return the existing file if it exists, otherwise it will return
// a brand new copy from the template.
func (r *ServerInterruptRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != ServerInterruptPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (r *ServerInterruptRender) IsFirst() bool {
	return r.prev == nil
}

// IsModified 代码是否更改
func (r *ServerInterruptRender) IsModified() bool {
	return false
}
