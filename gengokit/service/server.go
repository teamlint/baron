package service

import (
	"io"

	"github.com/pkg/errors"
	"github.com/teamlint/baron/gengokit"
)

func NewServer(prev io.Reader) gengokit.Renderable {
	return &ServerRender{
		prev: prev,
	}
}

// ServerRender 服务端
type ServerRender struct {
	prev io.Reader
}

func (r *ServerRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != ServerPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (r *ServerRender) IsFirst() bool {
	return r.prev == nil
}

// IsModified 代码是否更改
func (r *ServerRender) IsModified() bool {
	return false
}
