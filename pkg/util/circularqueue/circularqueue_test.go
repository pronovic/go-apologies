package circularqueue

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {
	var err error
	var entries []string

	queue := New[string](entries)

	err = queue.SetFirst("")
	assert.EqualError(t, err, "queue is empty")

	err = queue.SetFirst("a")
	assert.EqualError(t, err, "queue is empty")

	_, err = queue.Next()
	assert.EqualError(t, err, "queue is empty")
}

func TestSingle(t *testing.T) {
	var entry string
	var err error
	var entries = []string{"a"}

	queue := New[string](entries)

	for j := 0; j <= 10; j++ {
		entry, err = queue.Next()
		assert.Nil(t, err)
		assert.Equal(t, "a", entry)
	}

	err = queue.SetFirst("")
	assert.EqualError(t, err, "entry not found")

	err = queue.SetFirst("a")
	assert.Nil(t, err)
	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "a", entry)

	err = queue.SetFirst("d")
	assert.EqualError(t, err, "entry not found")

	err = queue.SetFirst("f")
	assert.EqualError(t, err, "entry not found")
}

func TestMultiple(t *testing.T) {
	var entry string
	var err error
	var entries = []string{"a", "b", "c", "d", "e"}

	queue := New[string](entries)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "a", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "b", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "c", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "d", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "e", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "a", entry)

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "b", entry)

	err = queue.SetFirst("")
	assert.EqualError(t, err, "entry not found")

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "c", entry) // this is the same place we would otherwise have been

	err = queue.SetFirst("a")
	assert.Nil(t, err)
	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "a", entry)
	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "b", entry)

	err = queue.SetFirst("d")
	assert.Nil(t, err)
	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "d", entry)
	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "e", entry)

	err = queue.SetFirst("f")
	assert.EqualError(t, err, "entry not found")

	entry, err = queue.Next()
	assert.Nil(t, err)
	assert.Equal(t, "a", entry) // this is the same place we would otherwise have been
}
