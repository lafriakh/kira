package log

type Fields map[string]any

// Get field value by name.
func (f Fields) Get(name string) any {
	return f[name]
}
