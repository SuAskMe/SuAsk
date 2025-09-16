package triemux

import (
	"errors"
	"strings"
)

type TrieMuxNode struct {
	children map[string]*TrieMuxNode
	isEnd    bool
}

type TrieMux struct {
	root *TrieMuxNode
}

func NewTrieMux() *TrieMux {
	return &TrieMux{
		root: &TrieMuxNode{
			children: make(map[string]*TrieMuxNode),
			isEnd:    false,
		},
	}
}

func (t *TrieMux) Insert(path string) error {
	if path[0] != '/' {
		return errors.New("path should start with '/'")
	}
	pathSp := strings.Split(path[1:], "/")
	node := t.root
	for _, word := range pathSp {
		if word == "" {
			continue
		}
		if _, ok := node.children[word]; !ok {
			node.children[word] = &TrieMuxNode{
				children: make(map[string]*TrieMuxNode),
				isEnd:    false,
			}
		}
		node = node.children[word]
	}
	node.isEnd = true
	return nil
}

func (t *TrieMux) getSplitIndexFrom(path *string, singleSep byte, st int) int {
	for i := st; i < len(*path); i++ {
		if (*path)[i] == singleSep {
			return i
		}
	}
	return len(*path)
}

func (t *TrieMux) HasPrefix(path string) bool {
	if path[0] != '/' {
		return false
	}
	node := t.root
	ed := 0
	for st := 1; st < len(path); st = ed + 1 {
		ed = t.getSplitIndexFrom(&path, '/', st)
		next, ok := node.children[path[st:ed]]
		// fmt.Println("next", next, "ok", ok, "path[st:ed]", path[st:ed])
		if next.isEnd {
			return true
		} else if !ok {
			return false
		}
		node = next
	}
	return false
}
