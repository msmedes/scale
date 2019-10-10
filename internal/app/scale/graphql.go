package scale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	pb "github.com/msmedes/scale/internal/app/scale/proto"
	"go.uber.org/zap"
)

type reqBody struct {
	Query string `json:"query"`
}

type nodeMetadata struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

type metadata struct {
	Node *nodeMetadata `json:"nodeMetadata"`
}

// GraphQL GraphQL object
type GraphQL struct {
	addr   string
	schema graphql.Schema
	logger *zap.SugaredLogger
	rpc    *RPC
}

// NewGraphQL Instantiate new GraphQL instance
func NewGraphQL(addr string, logger *zap.SugaredLogger, rpc *RPC) *GraphQL {
	obj := &GraphQL{
		addr:   addr,
		logger: logger,
		rpc:    rpc,
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

func (g *GraphQL) execute(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        g.schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		g.logger.Infof("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}

func (g *GraphQL) buildSchema() {
	nodeMetadataType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "NodeMetadata",
			Fields: graphql.Fields{
				"id":   &graphql.Field{Type: graphql.String},
				"addr": &graphql.Field{Type: graphql.String},
			},
		},
	)

	metadataType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Metadata",
			Fields: graphql.Fields{
				"node": &graphql.Field{
					Type: nodeMetadataType,
				},
			},
		},
	)

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"metadata": &graphql.Field{
					Type: metadataType,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						nodeMeta, err := g.rpc.GetNodeMetadata(context.Background(), &pb.Empty{})

						if err != nil {
							return nil, err
						}

						meta := &metadata{
							Node: &nodeMetadata{
								ID:   fmt.Sprintf("%x", nodeMeta.GetId()),
								Addr: nodeMeta.GetAddr(),
							},
						}

						return meta, nil
					},
				},
				"get": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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

	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"set": &graphql.Field{
					Type: graphql.Int,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
						"value": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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

	schema, _ := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	g.schema = schema
}
