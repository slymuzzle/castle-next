//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	ex, err := entgql.NewExtension(
		// Tell Ent to generate a GraphQL schema for
		// the Ent schema in a file named ent.graphql.
		entgql.WithSchemaGenerator(),
		entgql.WithWhereInputs(true),
		entgql.WithSchemaPath("graph/schema/ent.graphql"),
		entgql.WithConfigPath("gqlgen.yml"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	templates := entgql.AllTemplates
	templates = append(templates, gen.MustParse(
		gen.NewTemplate("pulid.tmpl").
			ParseFiles("./ent/schema/pulid/template/pulid.tmpl")),
	)

	opts := []entc.Option{
		entc.Extensions(ex),
	}

	if err := entc.Generate("./ent/schema", &gen.Config{
		Features: []gen.Feature{
			gen.FeatureIntercept,
			gen.FeatureSnapshot,
			gen.FeatureVersionedMigration,
			gen.FeatureUpsert,
		},
		Templates: templates,
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
