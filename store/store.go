package store

type Element struct {
	Value string
	Command string
}

type Store interface {
	Insert(key string, el Element) error
	Lookup(key string) (Element, error)
	Remove(key string) bool
}

type StoreFactory func() (*Store, error)
