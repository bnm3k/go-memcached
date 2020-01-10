package cache

//Cache ...
type Cache interface {
	Set(string, string, int)
	Get(string) (string, bool)
	Exists(string) bool
	Delete(string) error
}
