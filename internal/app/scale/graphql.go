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

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"get": &graphql.Field{
					Type: graphql.String,
					Args: graphql.FieldConfigArgument{
						"key": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						key := []byte(p.Args["key"].(string))

						res, err := rpc.GetLocal(context.Background(), &pb.GetRequest{Key: key})

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

						_, err := rpc.SetLocal(context.Background(), &pb.SetRequest{Key: key, Value: val})

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

	obj.schema = schema

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
