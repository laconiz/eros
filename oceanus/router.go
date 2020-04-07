package oceanus

type Key string
type Type string

type Router struct {
	Key  Key
	Type Type
}

type Routers struct {
	Keys map[Key]Router
}
