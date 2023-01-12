package main

import (
	"cloudinary_medias/api"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/cloudinary/cloudinary-go/v2"
)

func main() {
	cld, _ := cloudinary.NewFromParams("<your-cloud-name>", "<your-api-key>", "<your-api-secret>")

	db, err := gorm.Open(sqlite.Open("image.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	if err := db.AutoMigrate(&api.Image{}); err != nil {
		log.Fatalf("failed to migrate up")
	}

	apiServer := api.New(&api.RouterOptions{
		Cloudinary: cld,
		GormDB:     db,
	})

	if apiServer.Run(":8080"); err != nil {
		log.Fatalf("failed to run server")
	}
}
