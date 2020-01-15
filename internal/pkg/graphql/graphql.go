package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
	md "google.golang.org/grpc/metadata"

	gql "github.com/graphql-go/graphql"
	"github.com/msmedes/scale/internal/pkg/rpc"
	pb "github.com/msmedes/scale/internal/pkg/rpc/proto"
	"github.com/msmedes/scale/internal/pkg/trace"
	tracePb "github.com/msmedes/scale/internal/pkg/trace/proto"
	"go.uber.org/zap"
)

type reqBody struct {
	Query string `json:"query"`
}

// GraphQL GraphQL object
type GraphQL struct {
	addr        string
	schema      gql.Schema
	sugar       *zap.SugaredLogger
	logger      *zap.Logger
	rpc         *rpc.RPC
	traceClient tracePb.TraceClient
	traceConn   *grpc.ClientConn
	nodeAddr    string
}

// NewGraphQL Instantiate new GraphQL instance
func NewGraphQL(addr string, r *rpc.RPC, nodeAddr string) *GraphQL {
	obj := &GraphQL{addr: addr, rpc: r, nodeAddr: nodeAddr}

	obj.buildSchema()

	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	sugar := logger.Sugar()
	obj.sugar = sugar
	obj.logger = logger

	obj.traceClient, obj.traceConn = trace.NewClient("0.0.0.0:5000")

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

	g.sugar.Infof("listening: %s", g.addr)

	err := http.ListenAndServe(g.addr, nil)

	if err != nil {
		g.sugar.Fatalf("failed to listen: %v", err)
	}
}

func (g *GraphQL) execute(query string) *gql.Result {
	result := gql.Do(gql.Params{
		Schema:        g.schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		g.sugar.Infof("wrong result, unexpected errors: %v", result.Errors)
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
				"port":        &gql.Field{Type: gql.NewNonNull(gql.String)},
				"predecessor": &gql.Field{Type: remoteNodeMetadataType},
				"successor":   &gql.Field{Type: remoteNodeMetadataType},
				"fingerTable": &gql.Field{Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.String)))},
				"keys":        &gql.Field{Type: gql.NewNonNull(gql.NewList(gql.NewNonNull(gql.String)))},
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

						var ft []string

						for _, k := range nodeMeta.GetFingerTable() {
							ft = append(ft, fmt.Sprintf("%x", k))
						}

						node := &nodeMetadata{
							ID:          fmt.Sprintf("%x", nodeMeta.GetId()),
							Addr:        nodeMeta.GetAddr(),
							Port:        nodeMeta.GetPort(),
							FingerTable: ft,
							Keys:        nodeMeta.GetKeys(),
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

						id := "HELLO!"
						meta := md.Pairs("traceID", id)
						ctx := md.NewOutgoingContext(context.Background(), meta)
						g.traceClient.StartTrace(ctx, &tracePb.StartTraceRequest{TraceID: id})
						g.traceClient.AppendTrace(ctx, &tracePb.AppendTraceRequest{TraceID: id, Addr: g.nodeAddr})
						res, err := g.rpc.Get(ctx, &pb.GetRequest{Key: key})

						trace, traceErr := g.traceClient.GetTrace(context.Background(), &tracePb.TraceQuery{TraceID: id})
						g.sugar.Info(trace, traceErr)
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
						// id := "hello world"
						// ctx := md.AppendToOutgoingContext(context.Background(), "traceID", id)
						// g.traceClient.StartTrace(ctx, &tracePb.StartTraceRequest{TraceID: id})
						_, err := g.rpc.Set(context.Background(), &pb.SetRequest{Key: key, Value: val})

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

//Shutdown clean exit
func (g *GraphQL) Shutdown() {
	g.logger.Sync()
	g.sugar.Sync()

	g.traceConn.Close()
	g.sugar.Info("Closing connection to trace server")
}
