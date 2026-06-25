package main

import "fmt"

type Node struct {
	Key   string
	Value string
	Prev  *Node
	Next  *Node
}

type LRUCache struct {
	Capacity int
	Items    map[string]*Node
	Head     *Node
	Tail     *Node
}

func main() {
	cache := NewLRUCache(3)

	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3")

	fmt.Println(cache.Head.Key) // c
	fmt.Println(cache.Tail.Key) // a

	cache.Get("a")

	fmt.Println(cache.Head.Key) // a
	fmt.Println(cache.Tail.Key) // b

	cache.Set("d", "4")

	_, ok := cache.Get("b")
	fmt.Println(ok) // false
}

func NewLRUCache(capacity int) *LRUCache {
	cache := LRUCache{Capacity: capacity}
	cache.Items = make(map[string]*Node)

	return &cache
}

func (c *LRUCache) addToFront(node *Node) {
	node.Prev = nil
	node.Next = c.Head

	if c.Head != nil {
		c.Head.Prev = node
	} else {
		c.Tail = node
	}
	c.Head = node
}

func (c *LRUCache) removeNode(node *Node) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		c.Head = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		c.Tail = node.Prev
	}

	node.Prev = nil
	node.Next = nil
}

func (c *LRUCache) moveToFront(node *Node) {
	c.removeNode(node)
	c.addToFront(node)
}

func (c *LRUCache) Set(key string, value string) {
	node, ok := c.Items[key]
	if ok {
		node.Value = value
		c.moveToFront(node)
		return
	}

	node = &Node{
		Key:   key,
		Value: value,
	}

	c.addToFront(node)
	c.Items[key] = node

	if len(c.Items) > c.Capacity {
		oldTail := c.Tail
		c.removeNode(oldTail)
		delete(c.Items, oldTail.Key)
	}
}

func (c *LRUCache) Get(key string) (string, bool) {
	node, ok := c.Items[key]
	if !ok {
		return "", false
	}
	c.moveToFront(node)
	return node.Value, true
}
