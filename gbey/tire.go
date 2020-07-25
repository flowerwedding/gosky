package gbey
//实现Tire树的插入和查找功能，以运用到路由中去。
//之前的简单map存储路由表，使用map存储键值对，索引高效，但只能存储静态路由。
//动态路由就是一条路由规则可以匹配某一类型而非某一固定的路由。
//动态路由最常用的数据结构是前缀树（tire树），它的每一个节点的所有子路由都拥有相同的前缀。
//HTTP请求的路径刚好是由 / 分隔的多段构成，因此，每一段可以作为前缀树的一个节点。通过树结构查询，如果中间某一层的节点都不满足条件，就说明没有匹配到路由，查询结束。
import "strings"

//树节点上存储的信息
//为了实现动态路由匹配，在普通的树基础上加上了isWild参数。
//例如在匹配 /p/go/doc/ 这个路由时，第一层节点，p 精准匹配到了p，第二层节点，go 模糊匹配到 :lang，那么将会把lang这个参数赋值为go，继续匹配下一层。
type node struct {
	pattern  string // 待匹配路由，例如 /p/:lang  ，就是最开始的时候和方法封装在一起作为map的key的那个
	part     string // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，**用于插入**，用于插入就只要找第一个
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {//便利子节点，寻找第一个，路由后面那个部分就是它或者是可以匹配任何值的精准匹配
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于**查找**，所以返回了一堆，返回的数据类型也是结构体指针的切片
func (n *node) matchChildren(part string) []*node {
//new()和make()区别
//内置函数 new 分配空间。传递给new 函数的是一个类型，不是一个值。返回值是 指向这个新分配的零值的指针。
//内建函数 make 分配并且初始化 一个 slice, 或者 map 或者 chan 对象。 并且只能是这三种对象。 和 new 一样，第一个参数是 类型，不是一个值。 但是make 的返回值就是这个类型（即使一个引用类型），而不是指针。
	nodes := make([]*node, 0)
	for _, child := range n.children {//查询条件一样，就是一个返回了一个，一个返回了一堆
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//路由最重要的是注册和匹配。
//开发服务时，注册路由规则，映射handler；访问时，匹配路由规则，查找到对应的handler。
//tire树支持节点的插入与查询。
//插入功能：递归查找每一层的节点，如果没有匹配到当前的part节点就新建一个。在匹配结束时，可以使用n.pattern == ”“ 来判断路由规则是否匹配成功。
//插叙功能：同样递归查询每一层的节点，退出规则是，匹配到了 ' *' ，或者匹配到了第 len(parts) 层节点。

//插入功能：递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个。
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {//height初值为0，就是匹配到了极限的下一层
		n.pattern = pattern//然后就直接创建一个完整的路由pattern，parts是pattern的每一个一小部分，是用来查找用的，这里确定路由没找到，被匹配对象是 n 里面原来的值
		return
	}

	//新建待匹配路由pattern和路由的一部分part是两个不一样的，pattern直接建在n下，part找它的上层，上上层，建在那个下面。
	part := parts[height]//如果每一个路由的每一小段都一样就是找到了，part把每一小段都抽出来
	child := n.matchChild(part)//根据那一小段的名字去找它子节点上的值
	if child == nil {//没有匹配到档期那的part节点，就新建一个
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)//递归
}

//查询功能：递归查询每一层的节点，匹配到了'*'，或者匹配到了第len(parts)节点就退出
func (n *node) search(parts []string, height int) *node {
	// strings.HasPrefix() 用来检测字符串是否以指定的前缀开头。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//例如，/p/:lang/doc 的 p 和 :lang 节点的pattern都是空，所以可以使用n.pattern == "" 来判断路由规则是否匹配成功。
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