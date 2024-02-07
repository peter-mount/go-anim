package attributes

import (
	"sort"
	"sync"
)

type ImageAttributes interface {
	GetAttributeNames() []string
	GetAttribute(string) Attribute
	SetAttribute(Attribute)
	ForEach(f func(Attribute) error) error
}

type Attribute struct {
	Name string
	Type string
	Data []byte
}

type DefaultImageAttributes struct {
	mutex sync.Mutex
	attrs map[string]Attribute
}

func (m *DefaultImageAttributes) GetAttribute(n string) Attribute {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.attrs == nil {
		return Attribute{}
	}
	return m.attrs[n]
}

func (m *DefaultImageAttributes) SetAttribute(a Attribute) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.attrs == nil {
		m.attrs = make(map[string]Attribute)
	}
	m.attrs[a.Name] = a
}

func (m *DefaultImageAttributes) GetAttributeNames() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var names []string

	if m.attrs != nil {
		for k, _ := range m.attrs {
			names = append(names, k)
		}
	}

	sort.SliceStable(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	return names
}

func (m *DefaultImageAttributes) ForEach(f func(Attribute) error) error {
	for _, n := range m.GetAttributeNames() {
		if err := f(m.GetAttribute(n)); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultImageAttributes) Clone() *DefaultImageAttributes {
	n := &DefaultImageAttributes{}
	_ = m.ForEach(func(a Attribute) error {
		n.SetAttribute(a)
		return nil
	})
	return n
}
