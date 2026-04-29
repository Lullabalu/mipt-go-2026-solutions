//go:build !solution

package lrucache

import "container/list"

type pair struct {
	key   int
	value int
}

type LruCache struct {
	cap int
	lst *list.List
	mp  map[int]*list.Element
}

func New(cap int) Cache {
	lc := &LruCache{
		cap: cap,
		lst: list.New(),
		mp:  make(map[int]*list.Element),
	}

	return lc
}

func (lc *LruCache) Get(key int) (int, bool) {
	elem, ok := lc.mp[key]
	if !ok {
		return 0, false
	}
	lc.lst.MoveToBack(elem)
	pr := (*elem).Value.(*pair)
	return pr.value, true
}

func (lc *LruCache) Set(key, value int) {
	if lc.cap <= 0 {
		return
	}

	if elem, ok := lc.mp[key]; ok {
		pr := elem.Value.(*pair)
		pr.value = value
		lc.lst.MoveToBack(elem)
		return
	}

	if len(lc.mp) >= lc.cap {
		old := lc.lst.Front()
		if old != nil {
			delete(lc.mp, old.Value.(*pair).key)
			lc.lst.Remove(old)
		}
	}

	elem := lc.lst.PushBack(&pair{key: key, value: value})
	lc.mp[key] = elem
}

func (lc *LruCache) Clear() {
	lc.lst.Init()

	lc.mp = make(map[int]*list.Element)
}

func (lc *LruCache) Range(f func(key, value int) bool) {
	for el := lc.lst.Front(); el != nil; el = el.Next() {
		pr := el.Value.(*pair)
		if !f(pr.key, pr.value) {
			return
		}
	}
}
