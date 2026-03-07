package shell

import (
	"encoding/json"
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

// MkdirAll creates the directory and all missing parent directories.
func (fs *FS) MkdirAll(p string, hidden bool) error {
	p = path.Clean(p)
	if p == "/" {
		return nil
	}
	parts := strings.Split(strings.TrimPrefix(p, "/"), "/")
	cur := fs.root
	for _, part := range parts {
		child, ok := cur.children[part]
		if !ok {
			child = &Node{
				Name:     part,
				Type:     NodeDir,
				Hidden:   hidden,
				children: make(map[string]*Node),
			}
			cur.children[part] = child
		} else if child.Type != NodeDir {
			return fmt.Errorf("not a directory: %s", part)
		}
		cur = child
	}
	return nil
}

// WriteFile creates (or overwrites) a file at p, auto-creating parent dirs.
func (fs *FS) WriteFile(p, content string, hidden bool) error {
	p = path.Clean(p)
	parentPath := path.Dir(p)
	if err := fs.MkdirAll(parentPath, false); err != nil {
		return err
	}
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

// SerializedNode is used for JSON marshaling of the virtual FS.
type SerializedNode struct {
	Type        string `json:"type"`
	Content     string `json:"content,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`
	Permissions string `json:"permissions,omitempty"`
}

// Serialize walks the entire FS tree and returns a JSON snapshot.
// The result is a flat map of absolute path -> node attributes.
func (fs *FS) Serialize() (string, error) {
	nodes := map[string]SerializedNode{}
	var walk func(p string, node *Node)
	walk = func(p string, node *Node) {
		nodes[p] = SerializedNode{
			Type:        string(node.Type),
			Content:     node.Content,
			Hidden:      node.Hidden,
			Permissions: node.Permissions,
		}
		for name, child := range node.children {
			childPath := p + "/" + name
			if p == "/" {
				childPath = "/" + name
			}
			walk(childPath, child)
		}
	}
	// Walk all root children (not root itself)
	for name, child := range fs.root.children {
		walk("/"+name, child)
	}
	b, err := json.Marshal(nodes)
	return string(b), err
}

// DeserializeFS reconstructs a FS from a JSON snapshot produced by Serialize.
func DeserializeFS(data string) (*FS, error) {
	var nodes map[string]SerializedNode
	if err := json.Unmarshal([]byte(data), &nodes); err != nil {
		return nil, err
	}

	fs := NewFS()

	// Sort paths so parents are created before children
	paths := make([]string, 0, len(nodes))
	for p := range nodes {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		n := nodes[p]
		if n.Type == "dir" {
			if err := fs.Mkdir(p, n.Hidden); err != nil {
				return nil, fmt.Errorf("deserialize mkdir %s: %w", p, err)
			}
		} else {
			if err := fs.WriteFile(p, n.Content, n.Hidden); err != nil {
				return nil, fmt.Errorf("deserialize write %s: %w", p, err)
			}
			// Restore permissions if set
			if n.Permissions != "" {
				node, _ := fs.Stat(p)
				if node != nil {
					node.Permissions = n.Permissions
				}
			}
		}
	}
	return fs, nil
}
