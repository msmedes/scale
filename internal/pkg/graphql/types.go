package graphql

type remoteNodeMetadata struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type nodeMetadata struct {
	ID          string             `json:"id"`
	Addr        string             `json:"addr"`
	Predecessor remoteNodeMetadata `json:"predecessor"`
	Successor   remoteNodeMetadata `json:"successor"`
}

type metadata struct {
	Node *nodeMetadata `json:"nodeMetadata"`
}
