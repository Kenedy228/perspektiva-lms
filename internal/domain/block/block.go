package block

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/element"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Block struct {
	id       uuid.UUID
	title    string
	elements []*element.Element
}

func New(params Params) (*Block, error) {
	title := normalizeTitle(params.Title)

	if err := validateTitle(title); err != nil {
		return nil, err
	}

	if err := validateElements(params.Elements); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Block{
		id:       id,
		title:    title,
		elements: copyElements(params.Elements),
	}, nil
}

func (b *Block) ID() uuid.UUID {
	return b.id
}

func (b *Block) Title() string {
	return b.title
}

func (b *Block) Elements() []*element.Element {
	return copyElements(b.elements)
}

func (b *Block) InsertElementAt(pos int, elem *element.Element) error {
	if pos < 0 || pos > len(b.elements) {
		return fmt.Errorf("%w, детали: неверная позиция вставки", ErrInvalid)
	}

	if elem == nil {
		return fmt.Errorf("%w, детали: элемент для вставки должен существовать", ErrInvalid)
	}

	elements := slices.Insert(b.elements, pos, elem.Clone())
	if err := validateElementsDuplication(elements); err != nil {
		return err
	}

	if err := validateElementsLimit(elements, elementsLimit); err != nil {
		return err
	}

	b.elements = elements
	return nil
}

func (b *Block) RemoveElementAt(pos int) error {
	if pos < 0 || pos >= len(b.elements) {
		return fmt.Errorf("%w, детали: неверная позиция удаления", ErrInvalid)
	}

	b.elements = slices.Delete(b.elements, pos, pos+1)
	return nil
}

func (b *Block) UpdateElement(elem *element.Element) error {
	for i := range b.elements {
		if b.elements[i].ID() == elem.ID() {
			b.elements[i] = elem.Clone()
			return nil
		}
	}

	return fmt.Errorf("%w, детали: блок не содержит обновляемый элемент", ErrInvalid)
}

func (b *Block) MoveFromTo(from, to int) error {
	if from < 0 || from >= len(b.elements) {
		return fmt.Errorf("%w, детали: неверная начальная позиция перемещаемого элемента", ErrInvalid)
	}

	if to < 0 || to >= len(b.elements) {
		return fmt.Errorf("%w, детали: неверная конечная позиция перемещаемого элемента", ErrInvalid)
	}

	if from == to {
		return nil
	}

	el := b.elements[from]
	b.elements = slices.Delete(b.elements, from, from+1)
	b.elements = slices.Insert(b.elements, to, el)

	return nil
}

func (b *Block) Clone() *Block {
	return &Block{
		id:       b.id,
		title:    b.title,
		elements: copyElements(b.elements),
	}
}

func (b *Block) Copy() (*Block, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Block{
		id:       id,
		title:    b.title,
		elements: copyElements(b.elements),
	}, nil
}
