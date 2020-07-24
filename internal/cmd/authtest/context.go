package authtest

import (
	"context"
)

const (
	portKey = iota
	dbConnKey
)

func withPort(ctx context.Context, port string) context.Context {
	return context.WithValue(ctx, portKey, port)
}

func getPortFromCtx(ctx context.Context) string {
	return ctx.Value(portKey).(string)
}

func withDBConn(ctx context.Context, dbConn string) context.Context {
	return context.WithValue(ctx, dbConnKey, dbConn)
}

func getDBConnFromCtx(ctx context.Context) string {
	return ctx.Value(dbConnKey).(string)
}
