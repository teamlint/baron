package service

import (
	"io"

	"github.com/pkg/errors"
	"github.com/teamlint/baron/gengokit"
)

// NewHook returns a new HookRender
func NewHook(prev io.Reader) gengokit.Renderable {
	return &HookRender{
		prev: prev,
	}
}

type HookRender struct {
	prev io.Reader
}

// Render will return the existing file if it exists, otherwise it will return
// a brand new copy from the template.
func (r *HookRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != HookPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}
