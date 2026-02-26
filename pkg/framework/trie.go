package framework

import (
	"strings"
)

type node struct {
	pattern  string
	children []*node
	isWild   bool
	handler  any
}

func (n *node) insertChild(pattern string, isWild bool) *node {
	for _, child := range n.children {
		if child.pattern == pattern {
			return child
		}
	}
	child := &node{pattern: pattern, isWild: isWild}
	n.children = append(n.children, child)
	return child
}

func (n *node) insert(method, path string, handler any) {
	parts := parsePath(path)
	current := n
	for _, part := range parts {
		isWild := len(part) > 0 && (part[0] == ':' || part[0] == '*')
		current = current.insertChild(part, isWild)
	}
	current.handler = handler
}

func parsePath(path string) []string {
	vs := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
		}
	}
	return parts
}

func (n *node) search(parts []string) (*node, map[string]string) {
	params := make(map[string]string)
	current := n
	for _, part := range parts {
		var next *node
		for _, child := range current.children {
			if child.pattern == part || child.isWild {
				if child.isWild {
					params[child.pattern[1:]] = part
				}
				next = child
				break
			}
		}
		if next == nil {
			return nil, nil
		}
		current = next
	}
	if current.handler == nil {
		return nil, nil
	}
	return current, params
}
