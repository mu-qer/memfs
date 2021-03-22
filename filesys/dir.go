package filesys

import (
	"log"
	"os"
	"sync"

	"memfs/fuse"
	"memfs/fuse/fs"

	"golang.org/x/net/context"
)

type Dir struct {
	Node
	Files   map[string]*File
	SubDirs map[string]*Dir
	sync.Mutex
}

func (dir *Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Printf("called dir[%s]'s Attr()", dir.Name)
	attr.Inode = dir.Inode
	attr.Mode = os.ModeDir | 0644
	return nil
}

func (dir *Dir) Access(ctx context.Context, req *fuse.AccessRequest) error {
	log.Printf("called dir[%s]'s Access()", dir.Name)
	return nil
}

func (dir *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("called dir[%s]'s Lookup( %s )", dir.Name, name)
	return dir.lookup(name)
}

func (dir *Dir) lookup(name string) (fs.Node, error) {
	log.Printf("dir:%s|lookupname:%s", dir.Name, name)
	if f, ok := dir.Files[name]; ok {
		log.Printf("Found file:%s, inode:%d, size:%d", name, f.Inode, f.Size)
		return f, nil
	}

	if d, ok := dir.SubDirs[name]; ok {
		log.Printf("Found dir:%s, inode:%d, size:%d", name, d.Inode, d.Size)
		return d, nil
	}

	for _, subdir := range dir.SubDirs {
		fn, err := subdir.lookup(name)
		if err == nil {
			return fn, nil
		}
	}

	log.Printf("lookupname:%s not found.", name)
	return nil, fuse.ENOENT
}

//not necessory
func (dir *Dir) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Printf("called dir[%s]'s Open()", dir.Name)
	return dir, nil
}

func (dir *Dir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	log.Printf("called dir[%s]'s Mkdir()", dir.Name)
	if d, ok := dir.SubDirs[req.Name]; ok {
		log.Printf("dir[%s] has subdir:%s already.", dir.Name, req.Name)
		return d, nil
	}

	newDir := &Dir{
		Node:    Node{Name: req.Name, Inode: NewInode()},
		Files:   make(map[string]*File),
		SubDirs: make(map[string]*Dir),
	}
	dir.SubDirs[req.Name] = newDir
	return newDir, nil
}

func (dir *Dir) ReadDir(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("called dir[%s]'s ReadDir(), name:%s", dir.Name, name)
	return dir.lookup(name)
}

func (dir *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Printf("called dir[%s] ReadDirAll()", dir.Name)
	children := make([]fuse.Dirent, len(dir.Files)+len(dir.SubDirs))
	Idx := 0

	for _, f := range dir.Files {
		children[Idx] = fuse.Dirent{Inode: f.Inode, Type: fuse.DT_File, Name: f.Name}
		Idx += 1
	}

	for _, d := range dir.SubDirs {
		children[Idx] = fuse.Dirent{Inode: d.Inode, Type: fuse.DT_Dir, Name: d.Name}
		Idx += 1
	}

	log.Printf("ReadDirAll dir[%s], has [%d] objs", dir.Name, Idx)
	return children, nil
}

func (dir *Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	log.Printf("called dir[%s] Create(), create file name:%s", dir.Name, req.Name)
	if f, ok := dir.Files[req.Name]; ok {
		log.Printf("file[%s] already exists, forbide to create again.")
		return f, f, nil
	}

	f := NewFile(req.Name, 0755)
	dir.Files[req.Name] = f
	return f, f, nil
}

func (dir *Dir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	log.Printf("called dir[%s] Remove(), remove name:%s", dir.Name, req.Name)
	if _, ok := dir.Files[req.Name]; ok {
		log.Printf("remove file:%s", req.Name)
		delete(dir.Files, req.Name)
		return nil
	}

	if _, ok := dir.SubDirs[req.Name]; ok {
		log.Printf("remove dir:%s", req.Name)
		delete(dir.SubDirs, req.Name)
		return nil
	}

	log.Printf("not found dir/file named:%s", req.Name)
	return fuse.ENOENT
}

func (dir *Dir) Rename(ctx context.Context, req *fuse.RenameRequest, newDir fs.Node) error {
	log.Printf("called dir[%s] Rename(), oldname:%s, newname:%s, newDir:%+v", dir.Name, req.OldName, req.NewName, newDir)
	if _, ok := dir.Files[req.OldName]; ok {
		dir.Files[req.OldName].Name = req.NewName
		return nil
	}

	if _, ok := dir.SubDirs[req.OldName]; ok {
		dir.SubDirs[req.OldName].Name = req.NewName
		return nil
	}

	log.Printf("Rename, oldname:%s not found.", req.OldName)
	return fuse.ENOENT
}
