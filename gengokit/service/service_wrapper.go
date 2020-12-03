package service

import (
	"io"

	"github.com/pkg/errors"

	"github.com/teamlint/baron/gengokit"
)

// NewServiceWrapper returns a Renderable that renders the middlewares.go file.
func NewServiceWrapper(prev io.Reader) *ServiceWrapperRender {
	return &ServiceWrapperRender{prev}
}

// ServiceWrapper satisfies the gengokit.Renderable interface to render
// service wrapper.
type ServiceWrapperRender struct {
	prev io.Reader
}

// Render creates the wrapper.go file. With no previous version it renders
// the templates, if there was a previous version loaded in, it passes that
// through.
func (w *ServiceWrapperRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != ServiceWrapperPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if w.prev != nil {
		return w.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (w *ServiceWrapperRender) IsFirst() bool {
	return w.prev == nil
}

// IsModified 代码是否更改
func (m *ServiceWrapperRender) IsModified() bool {
	return false
}
