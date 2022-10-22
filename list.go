package gocache

import (
	"container/list"
	"time"
)

func (c *cache) LPush(k string, x any) {
	var obj *list.List
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		obj = list.New()
	} else {
		switch item.Object.(type) {
		case *list.List:
			obj = item.Object.(*list.List)
		default:
			obj = list.New()

		}
	}

	obj.PushBack(x)
	c.items[k] = Item{
		Object:     obj,
		Expiration: 0,
	}

	c.mu.Unlock()
}

func (c *cache) LPop(k string) (any, bool) {
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		c.mu.Unlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.Unlock()
			return nil, false
		}
	}
	switch item.Object.(type) {
	case *list.List:
		obj := item.Object.(*list.List)
		ele := obj.Back()
		obj.Remove(ele)
		item.Object = obj
		if obj.Len() == 0 {
			c.delete(k)
		} else {
			c.items[k] = item
		}
		c.mu.Unlock()
		return ele.Value, true
	default:
		c.mu.Unlock()
		return nil, false

	}
}

// Rpush
func (c *cache) RPush(k string, x any) {
	var obj *list.List
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		obj = list.New()
	} else {
		switch item.Object.(type) {
		case *list.List:
			obj = item.Object.(*list.List)
		default:
			obj = list.New()

		}
	}

	obj.PushFront(x)
	c.items[k] = Item{
		Object:     obj,
		Expiration: 0,
	}

	c.mu.Unlock()
}

func (c *cache) RPop(k string) (any, bool) {
	c.mu.Lock()
	item, found := c.items[k]
	if !found {
		c.mu.Unlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			c.mu.Unlock()
			return nil, false
		}
	}
	switch item.Object.(type) {
	case *list.List:
		obj := item.Object.(*list.List)
		ele := obj.Front()
		obj.Remove(ele)
		item.Object = obj
		if obj.Len() == 0 {
			c.delete(k)
		} else {
			c.items[k] = item
		}
		c.mu.Unlock()
		return ele.Value, true
	default:
		c.mu.Unlock()
		return nil, false

	}
}
