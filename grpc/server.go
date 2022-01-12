package grpc

import (
	pb "canvas-server/grpc/proto/go"
	"canvas-server/infra/cloud_storage"
	"context"
	"log"
	"net/http"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
)

type Server func(mux *http.ServeMux)

func NewServer(api pb.InternalAPIServer, authenticate Authenticate) Server {
	auth := func(server http.Handler) http.Handler {
		return applyMiddleware(
			server,
			authenticate)
	}

	interceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		res, err := handler(ctx, req)
		if err != nil {
			log.Printf("error %s", err.Error())
		}
		return res, err
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	pb.RegisterInternalAPIServer(server, api)
	reflection.Register(server)

	return func(mux *http.ServeMux) {
		mux.Handle("/", auth(grpcweb.WrapServer(
			server,
			grpcweb.WithOriginFunc(func(origin string) bool {
				return true
			}),
			grpcweb.WithAllowNonRootResource(true),
		)))
	}
}

func applyMiddleware(target http.Handler, handlers ...func(http.Handler) http.Handler) http.Handler {
	h := target
	for _, mw := range handlers {
		h = mw(h)
	}
	return h
}

type api struct {
	gcsClient cloud_storage.Client
}

func NewAPI(gcsClient cloud_storage.Client) pb.InternalAPIServer {
	return &api{
		gcsClient: gcsClient,
	}
}

func (a *api) SignedGsUrls(ctx context.Context, req *pb.SignedGsUrlsRequest) (*pb.SignedGsUrlsResponse, error) {
	results := make([]string, 0, len(req.GsUrls))

	for _, gsURL := range req.GsUrls {
		u, _ := url.Parse(gsURL)
		signedURL, _ := a.gcsClient.Signature(u)
		results = append(results, signedURL.String())
	}

	return &pb.SignedGsUrlsResponse{
		Urls: results,
	}, nil
}
