package filesys

import (
	"log"

	"memfs/fuse"
	"memfs/fuse/fs"
	"memfs/fuse/fuseutil"

	"golang.org/x/net/context"
)

type File struct {
	Node
	data []byte
}

func NewFile(name string, mode uint32) *File {
	newId := NewInode()

	return &File{
		Node: Node{Name: name, Inode: newId, Size: 0},
		data: make([]byte, 1024),
	}
}

func (f *File) Attr(ctx context.Context, attr *fuse.Attr) error {
	log.Printf("called file[%s] Attr().", f.Name)

	attr.Inode = f.Inode
	attr.Mode = 0755
	attr.Size = f.Size
	return nil
}

func (f *File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Printf("called file[%s] Read().", f.Name)
	fuseutil.HandleRead(req, resp, f.data)
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("called file[%s] ReadAll(), return data:%s, filesize:%d", f.Name, string(f.data), f.Size)
	return f.data, nil
}

func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Printf("called file[%s] Write(), data:%s", f.Name, string(req.Data))
	resp.Size = len(req.Data)
	oldSize := f.Size
	f.Size = uint64(resp.Size)

	if f.Size > oldSize {
		f.data = append(f.data, make([]byte, f.Size-oldSize)...)
	}
	copy(f.data, req.Data)
	f.data = f.data[:f.Size]

	log.Printf("Write() file[%s] data:%s", f.Name, string(f.data))
	return nil
}

func (f *File) Flush(ctx context.Context, req *fuse.FlushRequest) error {
	log.Printf("called file[%s] Flush.", f.Name)
	return nil
}

func (f *File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	log.Printf("called file[%s] Open(). return file:%+v", f.Name, f)
	return f, nil
}

func (f *File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	log.Printf("called file[%s] Release().", f.Name)
	return nil
}

func (f *File) Fsync(ctx context.Context, req *fuse.FsyncRequest) error {
	log.Printf("called file[%s] Release().", f.Name)
	return nil
}
