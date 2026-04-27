package collection_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/shared/collection"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("должен создать коллекцию с 0 элементов", func(t *testing.T) {
		//Arrange
		c := collection.NewOrderedClonable([]fixture{})

		//Assert
		assert.Empty(t, c.Items())
	})

	t.Run("должен создать коллекцию с 2 элементами, порядок после передачи фикстур не меняется", func(t *testing.T) {
		//Arrange
		fixtures := []fixture{
			fixture{id: uuid.New(), title: "first"},
			fixture{id: uuid.New(), title: "second"},
		}
		c := collection.NewOrderedClonable(fixtures)

		//Assert
		assert.Equal(t, len(c.Items()), 2)

		//fixture{id: uuid.New(), title: "first"},
		assert.Equal(t, c.Items()[0].title, "first")
		//fixture{id: uuid.New(), title: "second"},
		assert.Equal(t, c.Items()[1].title, "second")
	})
}

func TestInsertAt(t *testing.T) {
	t.Run("должен вернуть ошибку, если вставляем на позицию > len(items)", func(t *testing.T) {
		t.Run("элементы есть", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: uuid.New(), title: "new"}
			err := c.InsertAt(item, 100)

			//Assert
			assert.Error(t, err)
		})

		t.Run("элементов нет", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{})

			//Act
			item := fixture{id: uuid.New(), title: "new"}
			err := c.InsertAt(item, 500)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен вернуть ошибку, если вставляем на позицию < 0", func(t *testing.T) {
		t.Run("элементы есть", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: uuid.New(), title: "new"}
			err := c.InsertAt(item, -100)

			//Assert
			assert.Error(t, err)
		})

		t.Run("элементов нет", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{})

			//Act
			item := fixture{id: uuid.New(), title: "new"}
			err := c.InsertAt(item, -500)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен вставить элемент", func(t *testing.T) {
		t.Run("элементов нет, вставка в начало", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{})

			//Act
			item := fixture{id: uuid.New(), title: "item"}
			err := c.InsertAt(item, 0)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, c.Items()[0].title, "item")
		})

		t.Run("элементы есть, вставка в начало", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			})

			//Act
			item := fixture{id: uuid.New(), title: "item"}
			err := c.InsertAt(item, 0)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 3)

			// fixture{id: uuid.New(), title: "item"}
			assert.Equal(t, c.Items()[0].title, "item")

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[1].title, "first")

			// fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[2].title, "second")
		})

		t.Run("элементы есть, вставка в конец", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			})

			//Act
			item := fixture{id: uuid.New(), title: "item"}
			err := c.InsertAt(item, 2)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 3)

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[0].title, "first")

			// fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[1].title, "second")

			// fixture{id: uuid.New(), title: "item"}
			assert.Equal(t, c.Items()[2].title, "item")
		})

		t.Run("элементы есть, вставка в середину (элементы раздвигаются)", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			})

			//Act
			item := fixture{id: uuid.New(), title: "item"}
			err := c.InsertAt(item, 1)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 3)

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[0].title, "first")

			// fixture{id: uuid.New(), title: "item"}
			assert.Equal(t, c.Items()[1].title, "item")

			// fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[2].title, "second")
		})
	})
}

func TestRemove(t *testing.T) {
	t.Run("должен вернуть ошибку, если элемент для удаления не присутствует в контейнере", func(t *testing.T) {
		t.Run("элементы есть", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: uuid.New(), title: "remove"}
			err := c.Remove(item)

			//Assert
			assert.Error(t, err)
		})

		t.Run("элементов нет", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{})

			//Act
			item := fixture{id: uuid.New(), title: "remove"}
			err := c.Remove(item)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен удалить элемент, если элемент для удаления присутствует в контейнере", func(t *testing.T) {
		t.Run("элемент в начале", func(t *testing.T) {
			//Arrange
			removeID := uuid.New()
			fixtures := []fixture{
				fixture{id: removeID, title: "remove"},
				fixture{id: uuid.New(), title: "should stay"},
				fixture{id: uuid.New(), title: "second should stay"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: removeID, title: "different name, but should remove"}
			err := c.Remove(item)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 2)

			//fixture{id: uuid.New(), title: "should stay"},
			assert.Equal(t, c.Items()[0].title, "should stay")

			//fixture{id: uuid.New(), title: "second should stay"},
			assert.Equal(t, c.Items()[1].title, "second should stay")
		})

		t.Run("элемент в конце", func(t *testing.T) {
			//Arrange
			removeID := uuid.New()
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "should stay"},
				fixture{id: uuid.New(), title: "second should stay"},
				fixture{id: removeID, title: "remove"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: removeID, title: "different name, but should remove"}
			err := c.Remove(item)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 2)

			//fixture{id: uuid.New(), title: "should stay"},
			assert.Equal(t, c.Items()[0].title, "should stay")

			//fixture{id: uuid.New(), title: "second should stay"},
			assert.Equal(t, c.Items()[1].title, "second should stay")
		})

		t.Run("элемент в середине", func(t *testing.T) {
			//Arrange
			removeID := uuid.New()
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "should stay"},
				fixture{id: removeID, title: "remove"},
				fixture{id: uuid.New(), title: "second should stay"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: removeID, title: "different name, but should remove"}
			err := c.Remove(item)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, len(c.Items()), 2)

			//fixture{id: uuid.New(), title: "should stay"},
			assert.Equal(t, c.Items()[0].title, "should stay")

			//fixture{id: uuid.New(), title: "second should stay"},
			assert.Equal(t, c.Items()[1].title, "second should stay")
		})
	})
}

func TestUpdate(t *testing.T) {
	t.Run("должен вернуть ошибку, если элемент для обновления не присутствует в контейнере", func(t *testing.T) {
		t.Run("элементы есть", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: uuid.New(), title: "update"}
			err := c.Update(item)

			//Assert
			assert.Error(t, err)
		})

		t.Run("элементов нет", func(t *testing.T) {
			//Arrange
			c := collection.NewOrderedClonable([]fixture{})

			//Act
			item := fixture{id: uuid.New(), title: "update"}
			err := c.Update(item)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен обновить элемент, если элемент для обновления присутствует в контейнере", func(t *testing.T) {
		//Arrange
		updateID := uuid.New()
		fixtures := []fixture{
			fixture{id: uuid.New(), title: "should not change"},
			fixture{id: updateID, title: "update"},
			fixture{id: uuid.New(), title: "second should not change"},
		}
		c := collection.NewOrderedClonable(fixtures)

		//Act
		item := fixture{id: updateID, title: "updated name"}
		err := c.Update(item)

		//Assert
		assert.NoError(t, err)
		assert.Equal(t, len(c.Items()), 3)

		//fixture{id: uuid.New(), title: "should not change"},
		assert.Equal(t, c.Items()[0].title, "should not change")

		// fixture{id: updateID, title: "update"},
		assert.Equal(t, c.Items()[1].title, "updated name")

		//fixture{id: uuid.New(), title: "second should not change"},
		assert.Equal(t, c.Items()[2].title, "second should not change")
	})
}

func TestMove(t *testing.T) {
	t.Run("должен вернуть ошибку, если начальная позиция >= len(items)", func(t *testing.T) {
		t.Run("нет элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			errLen := c.Move(0, 0)
			errGtLen := c.Move(1, 0)

			//Assert
			assert.Error(t, errLen)
			assert.Error(t, errGtLen)
		})

		t.Run("есть элементы", func(t *testing.T) {
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			errLen := c.Move(2, 0)
			errGtLen := c.Move(3, 0)

			//Assert
			assert.Error(t, errLen)
			assert.Error(t, errGtLen)
		})
	})

	t.Run("должен вернуть ошибку, если начальная позиция < 0", func(t *testing.T) {
		t.Run("нет элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(-1, 0)

			//Assert
			assert.Error(t, err)
		})

		t.Run("есть элементы", func(t *testing.T) {
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(-1, 0)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен вернуть ошибку, если конечная позиция >= len(items)", func(t *testing.T) {
		t.Run("есть элементы", func(t *testing.T) {
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			errLen := c.Move(0, 2)
			errGtLen := c.Move(0, 3)

			//Assert
			assert.Error(t, errLen)
			assert.Error(t, errGtLen)
		})
	})

	t.Run("должен вернуть ошибку, если конечная позиция < 0", func(t *testing.T) {
		t.Run("есть элементы", func(t *testing.T) {
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(0, -5)

			//Assert
			assert.Error(t, err)
		})
	})

	t.Run("должен вернуть nil и не перемещать элементы, если границы корректные и равны между собой", func(t *testing.T) {
		fixtures := []fixture{
			fixture{id: uuid.New(), title: "first"},
			fixture{id: uuid.New(), title: "second"},
		}
		c := collection.NewOrderedClonable(fixtures)

		//Act
		err := c.Move(0, 0)

		//Assert
		assert.NoError(t, err)

		//fixture{id: uuid.New(), title: "first"},
		assert.Equal(t, c.Items()[0].title, "first")

		//fixture{id: uuid.New(), title: "second"},
		assert.Equal(t, c.Items()[1].title, "second")
	})

	t.Run("должен переместить элементы, если границы корректные", func(t *testing.T) {
		t.Run("нижняя граница < верхней границы, нечетное количество элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
				fixture{id: uuid.New(), title: "third"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(0, 2)

			//Assert
			assert.NoError(t, err)

			//fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[0].title, "second")

			//fixture{id: uuid.New(), title: "third"},
			assert.Equal(t, c.Items()[1].title, "third")

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[2].title, "first")
		})

		t.Run("нижняя граница > верхней границы, нечетное количество элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
				fixture{id: uuid.New(), title: "third"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(1, 0)

			//Assert
			assert.NoError(t, err)

			//fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[0].title, "second")

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[1].title, "first")

			//fixture{id: uuid.New(), title: "third"},
			assert.Equal(t, c.Items()[2].title, "third")
		})

		t.Run("нижняя граница < верхней границы, четное количество элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
				fixture{id: uuid.New(), title: "third"},
				fixture{id: uuid.New(), title: "fourth"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(1, 3)

			//Assert
			assert.NoError(t, err)

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[0].title, "first")

			//fixture{id: uuid.New(), title: "third"},
			assert.Equal(t, c.Items()[1].title, "third")

			//fixture{id: uuid.New(), title: "fourth"},
			assert.Equal(t, c.Items()[2].title, "fourth")

			//fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[3].title, "second")
		})

		t.Run("нижняя граница > верхней границы, нечетное количество элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
				fixture{id: uuid.New(), title: "third"},
				fixture{id: uuid.New(), title: "fourth"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			err := c.Move(2, 0)

			//Assert
			assert.NoError(t, err)

			//fixture{id: uuid.New(), title: "third"},
			assert.Equal(t, c.Items()[0].title, "third")

			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, c.Items()[1].title, "first")

			//fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, c.Items()[2].title, "second")

			//fixture{id: uuid.New(), title: "fourth"},
			assert.Equal(t, c.Items()[3].title, "fourth")
		})
	})
}

func TestItems(t *testing.T) {
	t.Run("возвращает коллекцию", func(t *testing.T) {
		t.Run("есть элементы", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "first"},
				fixture{id: uuid.New(), title: "second"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			items := c.Items()

			//Assert
			assert.NotEmpty(t, items)
			//fixture{id: uuid.New(), title: "first"},
			assert.Equal(t, items[0].title, "first")
			//fixture{id: uuid.New(), title: "second"},
			assert.Equal(t, items[1].title, "second")
		})

		t.Run("нет элементов", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			items := c.Items()

			//Assert
			assert.Empty(t, items)
		})
	})
}

func TestCount(t *testing.T) {
	t.Run("на пустой коллекции должен вернуть 0", func(t *testing.T) {
		//Arrange
		fixtures := []fixture{}
		c := collection.NewOrderedClonable(fixtures)

		//Act
		count := c.Count()

		//Assert
		assert.Equal(t, count, 0)
	})

	t.Run("на непустой коллекции должен вернуть количество элементов в срезе коллекции", func(t *testing.T) {
		//Arrange
		fixtures := []fixture{
			fixture{id: uuid.New(), title: "first"},
			fixture{id: uuid.New(), title: "second"},
		}
		c := collection.NewOrderedClonable(fixtures)

		//Act
		count := c.Count()

		//Assert
		assert.Equal(t, count, 2)
	})
}

func TestContains(t *testing.T) {
	t.Run("на пустой коллекции должен вернуть false", func(t *testing.T) {
		//Arrange
		fixtures := []fixture{}
		c := collection.NewOrderedClonable(fixtures)

		//Act
		item := fixture{id: uuid.New(), title: "contains"}
		has := c.Contains(item)

		//Assert
		assert.False(t, has)
	})

	t.Run("непустая коллекция", func(t *testing.T) {
		t.Run("элемент содержится в коллекции, должен вернуть true", func(t *testing.T) {
			//Arrange
			containsID := uuid.New()
			fixtures := []fixture{
				fixture{id: containsID, title: "contains"},
				fixture{id: uuid.New(), title: "item"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: containsID, title: "different title, but id the same"}
			has := c.Contains(item)

			//Assert
			assert.True(t, has)
		})

		t.Run("элемент не содержится в коллекции, должен вернуть false", func(t *testing.T) {
			//Arrange
			fixtures := []fixture{
				fixture{id: uuid.New(), title: "item"},
			}
			c := collection.NewOrderedClonable(fixtures)

			//Act
			item := fixture{id: uuid.New(), title: "not contains"}
			has := c.Contains(item)

			//Assert
			assert.False(t, has)
		})
	})
}
