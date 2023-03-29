package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Dynamic SQL

const dsn = "root:root@tcp(localhost:3306)/yoo?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/pkg/model",
	})
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect database: %v", err)
		return
	}

	g.UseDB(db)

	g.GenerateModelAs("users", "UserM")
	g.GenerateModelAs("templates", "TemplateM")

	g.Execute()
}
