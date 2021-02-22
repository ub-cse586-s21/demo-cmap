package cmap

import "testing"

func TestEmptyMap(t *testing.T) {
	cm := New()
	cm.Free()
}

const oneKey = "key"
const oneValue = "value"

func TestInsertOneItem(t *testing.T) {
	cm := New()
	_, ok := cm.Get(oneKey)
	if ok {
		t.Error("Map held key before insertion")
	}

	cm.Put(oneKey, oneValue)
	v, ok := cm.Get(oneKey)
	if !ok {
		t.Error("Map did not contain key after insertion")
	}
	if v != oneValue {
		t.Errorf("Map value was not as expected: %v vs %v", v, oneValue)
	}
}
