package graph

import (
	"context"
	mw "middleware"
	"model"

	ex "exception"

	"github.com/99designs/gqlgen/graphql"
	"github.com/rs/zerolog/log"
)

func IsAdmin(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	log.Info().Any("obj", obj).Msg("IsAdmin")
	user, ok := ctx.Value(mw.ContextKey("operator")).(*model.User)
	if !(ok && user.IsAdmin) {
		return nil, ex.ErrAuthFail
	}
	return next(ctx)
}
func IsLogin(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	log.Info().Any("obj", obj).Msg("IsLogIn")
	_, ok := ctx.Value(mw.ContextKey("operator")).(*model.User)
	if !ok {
		return nil, ex.ErrAuthFail
	}
	return next(ctx)
}
