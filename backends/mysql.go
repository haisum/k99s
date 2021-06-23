package backends

type MySQL struct {
	defaultBackend
}

func newMySQLBackend(name, namespace string) Backend {
	return MySQL{defaultBackend{
		Name:      name,
		Namespace: namespace,
	}}
}
