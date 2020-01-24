package graphql

type remoteNodeMetadata struct {
	ID   string `json:",omitempty"`
	Addr string `json:",omitempty"`
}

type nodeMetadata struct {
	ID               string              `json:",omitempty"`
	Addr             string              `json:",omitempty"`
	Port             string              `json:",omitempty"`
	Predecessor      *remoteNodeMetadata `json:",omitempty"`
	Successor        *remoteNodeMetadata `json:",omitempty"`
	FingerTableIDs   []string            `json:",omitempty"`
	FingerTableAddrs []string            `json:",omitempty"`
	Keys             []string            `json:",omitempty"`
}

type metadata struct {
	Node *nodeMetadata `json:",omitempty"`
}

type getRequest struct {
	Value string        `json:",omitempty"`
	Trace []*traceEntry `json:",omitempty"`
}

type setRequest struct {
	Count int           `json:",omitempty"`
	Trace []*traceEntry `json:",omitempty"`
}

type traceEntry struct {
	Addr         string `json:",omitempty"`
	FunctionCall string `json:",omitempty"`
	Duration     string `json:",omitempty"`
}

type network struct {
	Nodes []string `json:",omitempty"`
}
