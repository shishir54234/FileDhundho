package main

import (
	// "fmt"
	"errors"
	"fmt"
	
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// also need to have caching ?
type Engine struct {
	mu sync.Mutex

	
	root *Node
	current  *Node
};

type Node struct {
	parent *Node
	children []*Node
	metadata *NodeMetadata
	loaded bool
	err error 
};

type NodeMetadata struct {
	Name     string
	Path     string
	IsDir    bool
	Size     int64
	ModTime  time.Time
}


func NewNode(path string, parent *Node) (*Node, error){
	metadata, err:= NewNodeMetadata(path);
	if(err!=nil){
		fmt.Println(err);
		return nil, err;
	}
	nd:= &Node{
		parent: parent,
		children: []*Node{}, 
		metadata: metadata,
		loaded: false, 
		err: nil,
	}
	return nd, nil;

}
func NewNodeMetadata(path string) (*NodeMetadata, error){
	info, err:= os.Stat(path);
	if(err!=nil){
		return nil, err;
	}
	metadata:= &NodeMetadata{
		Name: path, 
		Path: path,
		IsDir: info.IsDir(),
		Size: info.Size(),
		ModTime: info.ModTime(),
		
	}
	return metadata, nil;
}
func NewEngine(path string) *Engine {

	rootNode, err:= NewNode(path, nil);
	if err!=nil {
		panic(err);
	}
	loadChildren(rootNode);
	return &Engine{
		root: rootNode,
		current: rootNode,
	};
}
func (e *Engine) ChangeDirectory(node *Node) {
	e.current = node;
}

func loadChildren(n *Node){
	if n.loaded || !n.metadata.IsDir {
		return
	}
	entries, err := os.ReadDir(n.metadata.Path);

	if err!=nil{
		n.err = err
		n.loaded = true
		return
	}
	for _, entry := range entries {
		fileName:=filepath.Join(n.metadata.Path, entry.Name()); 
		child, err:= NewNode(fileName, n);

		if(err!=nil){
			panic("The child failled for entry " + entry.Name() )
		}

		n.children = append(n.children, child)
	}

	n.loaded = true
}

func (e *Engine) List() ([]*Node, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	loadChildren(e.current)
	return e.current.children, e.current.err
}

func (e *Engine) Enter(idx int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	loadChildren(e.current)

	if idx < 0 || idx >= len(e.current.children) {
		return errors.New("index out of range")
	}

	n := e.current.children[idx]
	if !n.metadata.IsDir {
		return errors.New("not a directory")
	}

	e.current = n
	return nil
}

func (e *Engine) Up() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.current.parent == nil {
		return errors.New("already at root directory")
	}
	e.current = e.current.parent
	return nil
}

func(e *Engine) Search(query string)([]*Node, error){
	e.mu.Lock()
	defer e.mu.Unlock()
	loadChildren(e.current);
	var results []*Node
	// go through current directory and find matching files

	for _, child := range e.current.children {
			if  containsIgnoreCase(child.metadata.Name, query){
				results = append(results, child);
			}
		}
	return results, nil;
}


func containsIgnoreCase(str, substr string) bool {
	strLower := strings.ToLower(str)
	substrLower := strings.ToLower(substr)
	return strings.Contains(strLower, substrLower)
}



