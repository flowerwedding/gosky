package gosky

import "strings"

type node struct {
	pattern  string // 待匹配路由
	part     string // 路由中的一部分
	children []*node // 子节点
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入功能：新建待匹配路由pattern和路由的一部分part是两个不一样的，pattern直接建在n下，part递归找它的上层，上上层，建在那个下面。
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)//递归
}

//查询功能：递归查询每一层的节点，匹配到了'*'，或者匹配到了第len(parts)节点就退出
func (n *node) search(parts []string, height int) *node {
	// strings.HasPrefix() 用来检测字符串是否以指定的前缀开头。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)//递归
		if result != nil {
			return result
		}
	}

	return nil
}