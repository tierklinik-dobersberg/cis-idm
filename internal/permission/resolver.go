package permission


type Resolver interface {
	Resolve([]string) ([]string, error)
}

// Compile time checks
var (
	_ Resolver = (*NoTree)(nil)
	_ Resolver = (*Tree)(nil)
)