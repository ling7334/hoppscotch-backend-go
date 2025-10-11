package graph

import (
	"context"
	"log/slog"
	mw "middleware"
	"model"

	ex "exception"

	"github.com/99designs/gqlgen/graphql"
)

func IsAdmin(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	slog.Info("IsAdmin", "obj", obj)
	user, ok := ctx.Value(mw.ContextKey("operator")).(*model.User)
	if !(ok && user.IsAdmin) {
		return nil, ex.ErrAuthFail
	}
	return next(ctx)
}
func IsLogin(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	slog.Info("IsLogIn", "obj", obj)
	_, ok := ctx.Value(mw.ContextKey("operator")).(*model.User)
	if !ok {
		return nil, ex.ErrAuthFail
	}
	return next(ctx)
}
