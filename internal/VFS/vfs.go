package vfs

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type VFSNode struct {
	Name     string
	IsDir    bool
	Content  string
	Children map[string]*VFSNode
}

type VFS struct {
	Root     *VFSNode
	Current  *VFSNode
	Username string
}

func NewVFS() *VFS {
	root := &VFSNode{
		Name:     "/",
		IsDir:    true,
		Children: make(map[string]*VFSNode),
	}
	return &VFS{
		Root:     root,
		Current:  root,
		Username: "Default user",
	}
}

func (vfs *VFS) LoadFromCSV(filename string) error {
	if filename == "" {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 2 {
			continue
		}

		nodeType := record[0]
		path := record[1]
		content := ""
		if len(record) > 2 {
			content = record[2]
		}

		vfs.CreateNode(path, nodeType == "directory", content)
	}

	return nil
}

func (vfs *VFS) CreateNode(path string, isDir bool, content string) {
	parts := strings.Split(path, "/")
	var cleanParts []string
	for _, part := range parts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	current := vfs.Root
	for i, part := range cleanParts {
		isLast := i == len(cleanParts)-1

		if isLast {
			newNode := &VFSNode{
				Name:    part,
				IsDir:   isDir,
				Content: content,
			}
			if isDir {
				newNode.Children = make(map[string]*VFSNode)
			}
			current.Children[part] = newNode
		} else {
			if child, exists := current.Children[part]; exists {
				current = child
			} else {
				newDir := &VFSNode{
					Name:     part,
					IsDir:    true,
					Children: make(map[string]*VFSNode),
				}
				current.Children[part] = newDir
				current = newDir
			}
		}
	}
}

func (vfs *VFS) FindNode(path string) (*VFSNode, error) {
	if path == "" || path == "." {
		return vfs.Current, nil
	}
	if path == ".." {
		return vfs.Root, nil
	}
	if path == "/" {
		return vfs.Root, nil
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	parts := strings.Split(path, "/")
	var cleanParts []string
	for _, part := range parts {
		if part != "" {
			cleanParts = append(cleanParts, part)
		}
	}

	current := vfs.Current
	if strings.HasPrefix(path, "/") {
		current = vfs.Root
	}

	for _, part := range cleanParts {
		if part == ".." {
			current = vfs.Root
			continue
		}
		if child, exists := current.Children[part]; exists {
			current = child
		} else {
			return nil, fmt.Errorf("no such file or directory: %s", path)
		}
	}

	return current, nil
}

func (vfs *VFS) LS(path string) ([]string, error) {
	var target *VFSNode
	var err error

	if path == "" {
		target = vfs.Current
	} else {
		target, err = vfs.FindNode(path)
		if err != nil {
			return nil, err
		}
	}

	if !target.IsDir {
		return nil, fmt.Errorf("not a directory: %s", path)
	}

	var files []string
	for name := range target.Children {
		files = append(files, name)
	}
	return files, nil
}

func (vfs *VFS) CD(path string) error {
	node, err := vfs.FindNode(path)
	if err != nil {
		return err
	}
	if !node.IsDir {
		return fmt.Errorf("not a directory: %s", path)
	}
	vfs.Current = node
	return nil
}

func (vfs *VFS) GetContent(path string) (string, error) {
	node, err := vfs.FindNode(path)
	if err != nil {
		return "", err
	}
	if node.IsDir {
		return "", fmt.Errorf("is a directory: %s", path)
	}
	return node.Content, nil
}

func (vfs *VFS) Touch(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	vfs.Current.Children[filename] = &VFSNode{
		Name:    filename,
		IsDir:   false,
		Content: "",
	}
	return nil
}

func (vfs *VFS) GetCurrentPath() string {
	return vfs.Current.Name
}
