package service

import (
	"io"
	"strings"

	"github.com/teamlint/baron/gengokit"
	"github.com/teamlint/baron/gengokit/service/templates"
)

const HookPath = "service/hooks.gotemplate"

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
func (h *HookRender) Render(_ string, _ *gengokit.Data) (io.Reader, error) {
	if h.prev == nil {
		return strings.NewReader(templates.Hook), nil
	} else {
		return h.prev, nil
	}
}
