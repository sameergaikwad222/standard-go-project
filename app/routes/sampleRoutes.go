package routes

import (
	"sample-go-crud/app/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterSampleRoutes(rg *gin.RouterGroup, s *controllers.SampleController) {
	sampleRoutes := rg.Group("/samples")

	sampleRoutes.GET("", s.GetFilteredSamples)
	sampleRoutes.GET("/:id", s.GetOneSampleById)
	sampleRoutes.GET("/count", s.GetFilteredSampleCount)
	sampleRoutes.DELETE("/:id", s.DeleteSampleById)
	sampleRoutes.PATCH("/:id", s.UpdateSingleSamples)
	sampleRoutes.PATCH("", s.UpdateMultipleSamples)
	sampleRoutes.POST("", s.BulkInsertSamples)
}
