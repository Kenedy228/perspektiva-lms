package collection

import (
	"fmt"
	"slices"
)

type OrderedClonable[T Identifiable] struct {
	items []T
}

func NewOrderedClonable[T Identifiable](items []T) *OrderedClonable[T] {
	return &OrderedClonable[T]{
		items: items,
	}
}

func (c *OrderedClonable[T]) InsertAt(item T, pos int) error {
	if pos < 0 || pos > len(c.items) {
		return fmt.Errorf("неверная позиция вставки")
	}

	c.items = slices.Insert(c.items, pos, item)
	return nil
}

func (c *OrderedClonable[T]) Remove(item T) error {
	idx := slices.IndexFunc(c.items, func(current T) bool {
		return current.ID() == item.ID()
	})

	if idx == -1 {
		return fmt.Errorf("удаляемый элемент не существует в контейнере")
	}

	c.items = slices.Delete(c.items, idx, idx+1)
	return nil
}

func (c *OrderedClonable[T]) Update(item T) error {
	idx := slices.IndexFunc(c.items, func(current T) bool {
		return current.ID() == item.ID()
	})

	if idx == -1 {
		return fmt.Errorf("обновляемый элемент не существует в контейнере")
	}

	c.items[idx] = item
	return nil
}

func (c *OrderedClonable[T]) Move(from, to int) error {
	if from < 0 || from >= len(c.items) {
		return fmt.Errorf("неверная начальная позиция перемещаемого элемента")
	}

	if to < 0 || to >= len(c.items) {
		return fmt.Errorf("неверная конечная позиция перемещаемого элемента")
	}

	if from == to {
		return nil
	}

	item := c.items[from]
	c.items = slices.Delete(c.items, from, from+1)
	c.items = slices.Insert(c.items, to, item)

	return nil
}

func (c *OrderedClonable[T]) Contains(item T) bool {
	return slices.ContainsFunc(c.items, func(current T) bool {
		return current.ID() == item.ID()
	})
}

func (c *OrderedClonable[T]) Items() []T {
	return c.items
}

func (c *OrderedClonable[T]) Count() int {
	return len(c.items)
}
