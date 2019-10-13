package graphql

type remoteNodeMetadata struct {
	ID   string `json:",omitempty"`
	Addr string `json:",omitempty"`
}

type nodeMetadata struct {
	ID          string              `json:",omitempty"`
	Addr        string              `json:",omitempty"`
	Predecessor *remoteNodeMetadata `json:",omitempty"`
	Successor   *remoteNodeMetadata `json:",omitempty"`
	FingerTable []string            `json:",omitempty"`
}

type metadata struct {
	Node *nodeMetadata `json:",omitempty"`
}
