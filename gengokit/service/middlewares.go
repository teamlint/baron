package service

import (
	"io"

	"github.com/pkg/errors"

	"github.com/teamlint/baron/gengokit"
)

// NewMiddlewares returns a Renderable that renders the middlewares.go file.
func NewMiddlewares(prev io.Reader) *Middlewares {
	return &Middlewares{prev}
}

// Middlewares satisfies the gengokit.Renderable interface to render
// middlewares.
type Middlewares struct {
	prev io.Reader
}

// Render creates the middlewares.go file. With no previous version it renders
// the templates, if there was a previous version loaded in, it passes that
// through.
func (m *Middlewares) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != MiddlewaresPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if m.prev != nil {
		return m.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (m *Middlewares) IsFirst() bool {
	return m.prev == nil
}

// IsModified 代码是否更改
func (m *Middlewares) IsModified() bool {
	return false
}
