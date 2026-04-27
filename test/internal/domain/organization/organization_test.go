package organization_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/organization"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("empty name return error", func(t *testing.T) {
		//Arrange-Assert
		newOrganizationBuilder().withInn("1231231231").
			build(t, organization.ErrInvalid)
	})

	t.Run("invalid inn return error", func(t *testing.T) {
		//Arrange-Assert
		newOrganizationBuilder().withName("name").
			build(t, organization.ErrInvalid)
		newOrganizationBuilder().withInn("1234").
			withName("name").
			build(t, organization.ErrInvalid)
		newOrganizationBuilder().withInn("1234567890123458").
			withName("name").
			build(t, organization.ErrInvalid)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		org := newOrganizationBuilder().withInn("1234567890").
			withName("name").
			build(t, nil)

		//Assert
		assert.Equal(t, org.INN(), "1234567890")
		assert.Equal(t, org.Name(), "name")
		assert.False(t, org.CreatedAt().IsZero())
		assert.Equal(t, org.CreatedAt(), org.UpdatedAt())
	})
}

func TestChangeINN(t *testing.T) {
	t.Run("invalid inn should return err and not modify state", func(t *testing.T) {
		//Arrange
		org := newOrganizationBuilder().withInn("1234567890").
			withName("name").
			build(t, nil)
		oldUpdatedAt := org.UpdatedAt()

		//Act
		err := org.ChangeINN("123")

		//Assert
		assert.ErrorIs(t, err, organization.ErrInvalid)
		assert.Equal(t, org.INN(), "1234567890")
		assert.Equal(t, oldUpdatedAt, org.UpdatedAt())
		assert.Equal(t, org.CreatedAt(), org.UpdatedAt())
	})

	t.Run("valid inn should modify state", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			org := newOrganizationBuilder().withInn("1234567890").
				withName("name").
				build(t, nil)
			oldUpdatedAt := org.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := org.ChangeINN("2222222222")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, org.INN(), "2222222222")
			assert.True(t, oldUpdatedAt.Before(org.UpdatedAt()))
			assert.True(t, org.CreatedAt().Before(org.UpdatedAt()))
		})
	})
}

func TestRename(t *testing.T) {
	t.Run("invalid name should return err and not modify state", func(t *testing.T) {
		//Arrange
		org := newOrganizationBuilder().withInn("1234567890").
			withName("name").
			build(t, nil)
		oldUpdatedAt := org.UpdatedAt()

		//Act
		err := org.Rename("")

		//Assert
		assert.ErrorIs(t, err, organization.ErrInvalid)
		assert.Equal(t, org.Name(), "name")
		assert.Equal(t, oldUpdatedAt, org.UpdatedAt())
		assert.Equal(t, org.CreatedAt(), org.UpdatedAt())
	})

	t.Run("valid inn should modify state", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			org := newOrganizationBuilder().withInn("1234567890").
				withName("name").
				build(t, nil)
			oldUpdatedAt := org.UpdatedAt()

			//Act
			time.Sleep(time.Second * 10)
			err := org.Rename("new name")

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, org.Name(), "new name")
			assert.True(t, oldUpdatedAt.Before(org.UpdatedAt()))
			assert.True(t, org.CreatedAt().Before(org.UpdatedAt()))
		})
	})
}
