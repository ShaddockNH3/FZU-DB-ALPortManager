//go:build ignore

package main

import (
	"FZU-DB-ALPortManager/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithDefaultQuery | gen.WithoutContext,
	})

	g.ApplyBasic(
		model.ShipInfo{},
	)

	g.Execute()
}
