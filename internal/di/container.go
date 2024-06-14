package di

type Container interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	MustSet(key string, value interface{})
	MustGet(key string) interface{}

	Clear()
	Size() int
	IsEmpty() bool
}
