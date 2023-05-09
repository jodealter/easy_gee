package gee

import "strings"

// 节点，是url路径中的每个节点，比如 /hello/:kds/jode中的三个都是节点，并且又有父子关系
type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

// 用于插入的，插入的话，就插入第一个匹配的就好了，类似于如果第一二条线路都匹配，那么选择第一条线路就好了
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 用于查询的，查询出所有匹配的节点，然后从这些节点中选出一个可以使用的
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 新来的节点可以进行插入，没有就新建
func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 查询匹配到的节点，返回多个，但是只用第一个就行了
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
