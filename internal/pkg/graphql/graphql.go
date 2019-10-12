package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gql "github.com/graphql-go/graphql"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"go.uber.org/zap"
)

type reqBody struct {
	Query string `json:"query"`
}

// GraphQL GraphQL object
type GraphQL struct {
	addr   string
	schema gql.Schema
	logger *zap.SugaredLogger
	rpc    *rpc.RPC
}

// NewGraphQL Instantiate new GraphQL instance
func NewGraphQL(addr string, logger *zap.SugaredLogger, r *rpc.RPC) *GraphQL {
	obj := &GraphQL{
		addr:   addr,
		logger: logger,
		rpc:    r,
	}

	obj.buildSchema()

	return obj
}

// ServerListen start GraphQL server
func (g *GraphQL) ServerListen() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var t reqBody
		err := decoder.Decode(&t)

		if err != err {
			http.Error(w, "error parsing JSON request body", 400)
		}

		result := g.execute(t.Query)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	err := http.ListenAndServe(g.addr, nil)

	if err != nil {
		g.logger.Fatalf("failed to listen: %v", err)
	}
}

func (g *GraphQL) execute(query string) *gql.Result {
	result := gql.Do(gql.Params{
		Schema:        g.schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		g.logger.Infof("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}

func (g *GraphQL) buildSchema() {
	remoteNodeMetadataType := gql.NewObject(
		gql.ObjectConfig{
			Name: "RemoteNodeMetadata",
			Fields: gql.Fields{
				"id":   &gql.Field{Type: gql.NewNonNull(gql.String)},
				"addr": &gql.Field{Type: gql.NewNonNull(gql.String)},
			},
		},
	)

	nodeMetadataType := gql.NewObject(
		gql.ObjectConfig{
			Name: "NodeMetadata",
			Fields: gql.Fields{
				"id":          &gql.Field{Type: gql.NewNonNull(gql.String)},
				"addr":        &gql.Field{Type: gql.NewNonNull(gql.String)},
				"predecessor": &gql.Field{Type: remoteNodeMetadataType},
				"successor":   &gql.Field{Type: remoteNodeMetadataType},
			},
		},
	)

	metadataType := gql.NewObject(
		gql.ObjectConfig{
			Name: "Metadata",
			Fields: gql.Fields{
				"node": &gql.Field{
					Type: gql.NewNonNull(nodeMetadataType),
				},
			},
		},
	)

	queryType := gql.NewObject(
		gql.ObjectConfig{
			Name: "Query",
			Fields: gql.Fields{
				"metadata": &gql.Field{
					Type: metadataType,
					Resolve: func(p gql.ResolveParams) (interface{}, error) {
						nodeMeta, err := g.rpc.GetNodeMetadata(context.Background(), &pb.Empty{})

						node := &nodeMetadata{
							ID:   fmt.Sprintf("%x", nodeMeta.GetId()),
							Addr: nodeMeta.GetAddr(),
						}

						if err != nil {
							return nil, err
						}

						predID := nodeMeta.GetPredecessorId()
						succID := nodeMeta.GetSuccessorId()

						if predID != nil {
							node.Predecessor = &remoteNodeMetadata{
								ID:   fmt.Sprintf("%x", predID),
								Addr: nodeMeta.GetPredecessorAddr(),
							}
						}

						if succID != nil {
							node.Successor = &remoteNodeMetadata{
								ID:   fmt.Sprintf("%x", succID),
								Addr: nodeMeta.GetSuccessorAddr(),
							}
						}

						meta := &metadata{Node: node}

						return meta, nil
					},
				},
				"get": &gql.Field{
					Type: gql.String,
					Args: gql.FieldConfigArgument{
						"key": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
					},
					Resolve: func(p gql.ResolveParams) (interface{}, error) {
						key := []byte(p.Args["key"].(string))

						res, err := g.rpc.GetLocal(context.Background(), &pb.GetRequest{Key: key})

						if err != nil {
							return nil, err
						}

						return fmt.Sprintf("%s", res.Value), nil
					},
				},
			},
		},
	)

	mutationType := gql.NewObject(
		gql.ObjectConfig{
			Name: "Mutation",
			Fields: gql.Fields{
				"set": &gql.Field{
					Type: gql.Int,
					Args: gql.FieldConfigArgument{
						"key":   &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
						"value": &gql.ArgumentConfig{Type: gql.NewNonNull(gql.String)},
					},
					Resolve: func(p gql.ResolveParams) (interface{}, error) {
						key := []byte(p.Args["key"].(string))
						val := []byte(p.Args["value"].(string))

						_, err := g.rpc.SetLocal(context.Background(), &pb.SetRequest{Key: key, Value: val})

						if err != nil {
							return nil, err
						}

						return 1, nil
					},
				},
			},
		},
	)

	schema, _ := gql.NewSchema(
		gql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	g.schema = schema
}
