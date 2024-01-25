package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func GetPreloadedDB(db *gorm.DB, ctx context.Context) *gorm.DB {
	for _, c := range GetPreloads(ctx) {
		db = db.Preload(c)
	}
	return db
}

func GetPreloads(ctx context.Context) []string {
	return GetNestedPreloads(
		graphql.GetOperationContext(ctx),
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, cases.Title(language.English).String(column.Name))
		children := graphql.CollectFields(ctx, column.Selections, nil)
		if len(children) > 0 {
			preloads = append(preloads, prefixColumn)
			preloads = append(preloads, GetNestedPreloads(ctx, children, prefixColumn)...)
		}
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
