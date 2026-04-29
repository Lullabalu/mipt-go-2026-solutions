//go:build !solution

package treeiter

type Node[T any] interface {
	comparable
	Left() T
	Right() T
}

func Dfs[T Node[T]](cur T, f func(node T)) {
	var zero T
	if zero == cur {
		return
	}

	Dfs(cur.Left(), f)
	f(cur)
	Dfs(cur.Right(), f)
}

func DoInOrder[T Node[T]](root T, f func(node T)) {
	Dfs(root, f)
}
