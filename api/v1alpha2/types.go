package v1alpha2

// NamespaceService is an interface for management of namespaces.
type NamespaceService interface {
	List(prefix string) (NamespaceList, error)
	Get(name string) (Namespace, error)
	Set(namespace *Namespace) (Namespace, error)
	Delete(name string) (Namespace, error)
}

// NamespaceList is an array/list of Namespaces.
type NamespaceList []Namespace

// Namespace contains the definition of namespace.
type Namespace struct {
	Name string `json:"name"`
}

// RequestError contains an error message returned to user.
type RequestError struct {
	Error string `json:"error"`
}
