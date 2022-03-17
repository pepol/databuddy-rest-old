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
	Name   string          `json:"name"`
	Spec   NamespaceSpec   `json:"spec"`
	Status NamespaceStatus `json:"status"`
}

// NamespaceSpec contains the specification of namespace object.
type NamespaceSpec struct{}

// NamespaceStatus contains the status of namespace object.
type NamespaceStatus struct {
	Collections CollectionList `json:"collections"`
}

// CollectionService is an interface for management of collections.
type CollectionService interface {
	List(namespace, prefix string) (CollectionList, error)
	Get(namespace, name string) (Collection, error)
	Set(collection *Collection) (Collection, error)
	Delete(namespace, name string) (Collection, error)
}

// CollectionList is an array/list of Collections.
type CollectionList []Collection

// Collection contains the definition of collection.
type Collection struct {
	Name      string         `json:"name"`
	Namespace string         `json:"namespace"`
	Spec      CollectionSpec `json:"spec"`
}

// CollectionSpec contains the specification of collection object.
type CollectionSpec struct{}

// RequestError contains an error message returned to user.
type RequestError struct {
	Error string `json:"error"`
}
