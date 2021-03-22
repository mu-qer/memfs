package filesys

import (
	"memfs/fuse/fs"
)

type MemFS struct {
	root *Dir
}

var _ fs.FS = (*MemFS)(nil)

func (m *MemFS) Root() (fs.Node, error) {
	return m.root, nil
}

func NewMemFS(path string) *MemFS {
	return &MemFS{
		root: &Dir{Node: Node{Name: path, Inode: NewInode()},
			Files:   make(map[string]*File),
			SubDirs: make(map[string]*Dir),
		},
	}
}
