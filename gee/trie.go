package gee

import "strings"

// /:lang/doc
// /:lang/tutorial
// /:lang/intro
type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
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

// pattern	hello/world
// parts	[hello,world]
// height	0
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]       // height = 0, part = hello
	child := n.matchChild(part) // n = &node{pattern: "", part: "", children: nil, isWild: false}
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'} // child = &node{pattern: "", part: "hello", children: nil, isWild: false}
		n.children = append(n.children, child)                              // n = &node{pattern: "", part: "hello", children: &node{pattern: "", part: "hello", children: nil, isWild: false}, isWild: false}
	}
	child.insert(pattern, parts, height+1)
}

// 第一轮
// height = 0, part = hello
// n = &node{pattern: "", part: "", children: nil, isWild: false}
// child = &node{pattern: "", part: "hello", children: nil, isWild: false}
// n = &node{pattern: "", part: "", children: &node{pattern: "", part: "hello", children: nil, isWild: false}, isWild: false}

// 第二轮
// height = 1, part = world
// n = &node{pattern: "", part: "hello", children: nil, isWild: false} 上一步的child
// child = &node{pattern: "", part: "world", children: nil, isWild: false}
// n = &node{pattern: "", part: "hello", children: &node{pattern: "", part: "world", children: nil, isWild: false}, isWild: false}

/*
&node{
	pattern: "",
	part:    "",
	children: []*node{
		{
			pattern: "",
			part:    "hello",
			children: []*node{
				{
					pattern:  "/hello/world",
					part:     "world",
					children: nil,
					isWild:   false,
				},
				{
					pattern:  "/hello/world2",
					part:     "world2",
					children: nil,
					isWild:   false,
				},
			},
			isWild: false,
		},
		{
			pattern: "",
			part:    "p",
			children: []*node{
				{
					pattern:  "/p/doc",
					part:     "doc",
					children: nil,
					isWild:   false,
				},
			},
			isWild: false,
		},
	},
	isWild: false,
}
*/

// parts	[hello,world]
// height	0
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height] // height = 0, part = hello
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// height = 0, part = hello
// height = 1, part = world
