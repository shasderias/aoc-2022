package fs

import (
	"fmt"
	"path"
)

var ErrNotFound = fmt.Errorf("not found")

func New() *Context {
	root := newRoot()
	return &Context{
		root: root,
		cwd:  root,
	}
}

func newRoot() *Root {
	return &Root{
		children: make([]Node, 0),
	}
}

type Root struct {
	children []Node
}

func (r *Root) Name() string     { return "/" }
func (r *Root) Path() string     { return "/" }
func (r *Root) Parent() DirNode  { return nil }
func (r *Root) Children() []Node { return r.children }

func (r *Root) Create(n Node) error {
	if t := r.Open(n.Name()); t != nil {
		return fmt.Errorf("directory or file with name %s already exists", n.Name())
	}
	r.children = append(r.children, n)
	return nil
}

func (r *Root) Open(name string) Node {
	for _, n := range r.children {
		if n.Name() == name {
			return n
		}
	}
	return nil
}

func (r *Root) Size() int {
	totalSize := 0
	for _, n := range r.children {
		totalSize += n.Size()
	}
	return totalSize
}

type Dir struct {
	name     string
	parent   DirNode
	children []Node
}

func newDir(name string, parent DirNode) *Dir {
	return &Dir{
		name:     name,
		parent:   parent,
		children: make([]Node, 0),
	}
}

func (d *Dir) Name() string     { return d.name }
func (d *Dir) Path() string     { return path.Join(d.parent.Path(), d.name) }
func (d *Dir) Parent() DirNode  { return d.parent }
func (d *Dir) Children() []Node { return d.children }

func (d *Dir) Size() int {
	totalSize := 0
	for _, n := range d.children {
		totalSize += n.Size()
	}
	return totalSize
}

func (d *Dir) Create(n Node) error {
	d.children = append(d.children, n)
	return nil
}

func (d *Dir) Open(name string) Node {
	for _, n := range d.children {
		if n.Name() == name {
			return n
		}
	}
	return nil
}

func newFile(name string, size int) *File {
	return &File{name: name, size: size}
}

type File struct {
	name string
	size int
}

func (f *File) Name() string { return f.name }
func (f *File) Size() int    { return f.size }

type Node interface {
	Name() string
	Size() int
}

type DirNode interface {
	Node
	Path() string
	Parent() DirNode
	Children() []Node
	Create(n Node) error
	Open(name string) Node
}

type Context struct {
	root *Root
	cwd  DirNode
}

func (c *Context) CD(path string) error {
	switch path {
	case "/":
		c.cwd = c.root
	case "..":
		if c.cwd.Parent() == nil {
			return fmt.Errorf("cannot cd up from root")
		}
		c.cwd = c.cwd.Parent()
	default:
		target := c.cwd.Open(path)
		switch t := target.(type) {
		case DirNode:
			c.cwd = t
		case *File:
			return fmt.Errorf("cannot cd into file")
		case nil:
			return ErrNotFound
		}
	}
	return nil
}

func (c *Context) MkDir(name string) error {
	d := newDir(name, c.cwd)
	return c.cwd.Create(d)
}

func (c *Context) MkFile(name string, size int) error {
	f := newFile(name, size)
	return c.cwd.Create(f)
}

func (c *Context) CWD() DirNode {
	return c.cwd
}

type walkFunc func(n Node, path string) error

func (c *Context) Walk(fn walkFunc) {
	fn(c.cwd, c.cwd.Path())
	walk(c.cwd, fn, 0)
}

func walk(start DirNode, fn walkFunc, depth int) {
	for _, n := range start.Children() {
		switch t := n.(type) {
		case DirNode:
			fn(t, t.Path())
			walk(t, fn, depth+1)
		case *File:
			fn(t, path.Join(start.Path(), t.Name()))
			continue
		default:
			panic(fmt.Sprintf("unexpected node type %T", n))
		}
	}
}
