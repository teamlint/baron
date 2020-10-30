package service

import (
	"io"

	"github.com/pkg/errors"

	"github.com/teamlint/baron/gengokit"
)

// NewCmdClient 生成客户端CLI
func NewCmdClient(prev io.Reader) gengokit.Renderable {
	return &CmdClientRender{
		prev: prev,
	}
}

// CmdClientRender 客户端CLI
type CmdClientRender struct {
	prev io.Reader
}

// Render 生成代码, 实现 gengokit.Renderable 接口
func (r *CmdClientRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != CmdClientPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}
