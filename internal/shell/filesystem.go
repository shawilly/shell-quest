package shell

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

type NodeType string

const (
	NodeDir  NodeType = "dir"
	NodeFile NodeType = "file"
)

type Node struct {
	Name        string
	Type        NodeType
	Content     string
	Hidden      bool
	Permissions string
	children    map[string]*Node
}

type FS struct {
	root *Node
}

func NewFS() *FS {
	return &FS{
		root: &Node{
			Name:     "/",
			Type:     NodeDir,
			children: make(map[string]*Node),
		},
	}
}

func (fs *FS) resolve(p string) (*Node, error) {
	p = path.Clean(p)
	if p == "/" {
		return fs.root, nil
	}
	parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
	cur := fs.root
	for _, part := range parts {
		if cur.Type != NodeDir {
			return nil, fmt.Errorf("not a directory: %s", cur.Name)
		}
		child, ok := cur.children[part]
		if !ok {
			return nil, fmt.Errorf("no such file or directory: %s", p)
		}
		cur = child
	}
	return cur, nil
}

func (fs *FS) resolveParent(p string) (*Node, string, error) {
	p = path.Clean(p)
	parent := path.Dir(p)
	name := path.Base(p)
	node, err := fs.resolve(parent)
	if err != nil {
		return nil, "", fmt.Errorf("parent directory not found: %s", parent)
	}
	if node.Type != NodeDir {
		return nil, "", fmt.Errorf("not a directory: %s", parent)
	}
	return node, name, nil
}

func (fs *FS) Stat(p string) (*Node, error) {
	return fs.resolve(p)
}

func (fs *FS) Mkdir(p string, hidden bool) error {
	parent, name, err := fs.resolveParent(p)
	if err != nil {
		return err
	}
	parent.children[name] = &Node{
		Name:     name,
		Type:     NodeDir,
		Hidden:   hidden,
		children: make(map[string]*Node),
	}
	return nil
}

func (fs *FS) WriteFile(p, content string, hidden bool) error {
	parent, name, err := fs.resolveParent(p)
	if err != nil {
		return err
	}
	parent.children[name] = &Node{
		Name:    name,
		Type:    NodeFile,
		Content: content,
		Hidden:  hidden,
	}
	return nil
}

func (fs *FS) listDir(p string, showHidden bool) ([]*Node, error) {
	node, err := fs.resolve(p)
	if err != nil {
		return nil, err
	}
	if node.Type != NodeDir {
		return nil, fmt.Errorf("not a directory: %s", p)
	}
	var entries []*Node
	for _, child := range node.children {
		if !showHidden && child.Hidden {
			continue
		}
		entries = append(entries, child)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})
	return entries, nil
}

func (fs *FS) ListDir(p string) ([]*Node, error)    { return fs.listDir(p, false) }
func (fs *FS) ListDirAll(p string) ([]*Node, error) { return fs.listDir(p, true) }

// ResolvePath resolves p relative to cwd. Absolute paths are cleaned as-is.
func ResolvePath(cwd, p string) string {
	if strings.HasPrefix(p, "/") {
		return path.Clean(p)
	}
	return path.Clean(cwd + "/" + p)
}

func (fs *FS) Remove(p string) error {
	parent, name, err := fs.resolveParent(p)
	if err != nil {
		return err
	}
	if _, ok := parent.children[name]; !ok {
		return fmt.Errorf("no such file or directory: %s", p)
	}
	delete(parent.children, name)
	return nil
}

func (fs *FS) Copy(src, dst string) error {
	srcNode, err := fs.resolve(src)
	if err != nil {
		return err
	}
	if srcNode.Type == NodeDir {
		return fmt.Errorf("cp: cannot copy directory (use -r): %s", src)
	}
	return fs.WriteFile(dst, srcNode.Content, srcNode.Hidden)
}

func (fs *FS) Move(src, dst string) error {
	srcNode, err := fs.resolve(src)
	if err != nil {
		return err
	}
	dstParent, dstName, err := fs.resolveParent(dst)
	if err != nil {
		return err
	}
	dstParent.children[dstName] = srcNode
	srcNode.Name = dstName
	return fs.Remove(src)
}
