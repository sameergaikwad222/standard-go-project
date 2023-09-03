package inject

import (
	"context"
	"sample-go-crud/app/controllers"
	"sample-go-crud/app/repositories"
	"sample-go-crud/app/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type DI struct {
	Db *mongo.Database
}

func InitDI(db *mongo.Database) *DI {
	return &DI{
		Db: db,
	}
}

func (DI *DI) Inject(ctx context.Context, router *gin.Engine) {
	basePath := router.Group("/api/v1")
	controllers.InitSampleCache()
	sampleRepo := repositories.NewSampleRepository(DI.Db)
	sampleController := controllers.InitSampleController(sampleRepo, ctx)
	routes.RegisterSampleRoutes(basePath, sampleController)
}
