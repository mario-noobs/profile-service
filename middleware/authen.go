package middleware

import (
	"context"
	"demo-service/common"
	"github.com/gin-gonic/gin"
	"github.com/viettranx/service-context/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"strings"
)

type AuthClient interface {
	IntrospectToken(ctx context.Context, accessToken string) (sub string, tid string, err error)
}

func AuthInterceptor(ac AuthClient) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		if info.FullMethod == "/pb.UserService/CreateUser" {
			// Skip the auth interceptor for this method
			return handler(ctx, req)
		}
		// Extract metadata (headers)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, grpc.Errorf(codes.Unauthenticated, "missing metadata")
		}

		// Get the Authorization header
		token := md["authorization"]
		if len(token) == 0 {
			return nil, grpc.Errorf(codes.Unauthenticated, "missing authorization token")
		}

		sub, tid, err := ac.IntrospectToken(ctx, token[0])
		if err != nil {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		requester := core.NewRequester(sub, tid)

		ctx = context.WithValue(ctx, core.KeyRequester, requester)

		// Call the handler
		return handler(ctx, req)
	}
}

func RequireAuth(ac AuthClient) func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			common.WriteErrorResponse(c, err)
			c.Abort()
			return
		}

		sub, tid, err := ac.IntrospectToken(c.Request.Context(), token)

		if err != nil {
			common.WriteErrorResponse(c, core.ErrUnauthorized.WithDebug(err.Error()))
			c.Abort()
			return
		}

		c.Set(core.KeyRequester, core.NewRequester(sub, tid))

		c.Next()
	}
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"

	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" || strings.TrimSpace(parts[1]) == "null" {
		return "", core.ErrUnauthorized.WithError("missing access token")
	}

	return parts[1], nil
}
