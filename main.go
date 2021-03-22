package main

import (
	"fmt"
	"log"
	"os"

	"memfs/filesys"

	"memfs/fuse"
	"memfs/fuse/fs"
)

var inodeCount uint64

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Provide directory for mounting")
		os.Exit(1)
	}

	mountpoint := os.Args[1]

	conn, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	server := fs.New(conn, nil)
	fsys := filesys.NewMemFS(mountpoint)
	log.Println("About to serve memfs.")
	if err := server.Serve(fsys); err != nil {
		log.Panicln(err)
	}

	<-conn.Ready
	if err := conn.MountError; err != nil {
		log.Panicln(err)
	}
}
