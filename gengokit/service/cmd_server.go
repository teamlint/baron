package service

import (
	"io"

	"github.com/pkg/errors"

	"github.com/teamlint/baron/gengokit"
)

// NewCmdServer 生成服务端CLI
func NewCmdServer(prev io.Reader) gengokit.Renderable {
	return &CmdServerRender{
		prev: prev,
	}
}

// CmdServerRender 服务端CLI
type CmdServerRender struct {
	prev io.Reader
}

// Render 生成代码, 实现 gengokit.Renderable 接口
func (r *CmdServerRender) Render(path string, data *gengokit.Data) (io.Reader, error) {
	if path != CmdServerPath {
		return nil, errors.Errorf("cannot render unknown file: %q", path)
	}
	if r.prev != nil {
		return r.prev, nil
	}
	return data.ApplyTemplateFromPath(path)
}

// IsFirst 首次生成
func (r *CmdServerRender) IsFirst() bool {
	return r.prev == nil
}

// IsModified 代码是否更改
func (r *CmdServerRender) IsModified() bool {
	return false
}
