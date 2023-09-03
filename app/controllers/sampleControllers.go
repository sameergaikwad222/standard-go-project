package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sample-go-crud/app/helpers"
	"sample-go-crud/app/repositories"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
)

type SampleController struct {
	SampleRepository *repositories.SampleRepository
	Context          context.Context
}

var sampleCache *cache.Cache

func InitSampleController(samplerepository *repositories.SampleRepository, ctx context.Context) *SampleController {
	return &SampleController{
		SampleRepository: samplerepository,
		Context:          ctx,
	}
}

func InitSampleCache() {
	sampleCache = cache.New(60*time.Minute, 120*time.Minute)
}

func GetSampleCache(key string) (interface{}, bool) {
	return sampleCache.Get(key)
}

func SetSampleCache(key string, value interface{}) {
	sampleCache.Set(key, value, cache.DefaultExpiration)
}

func (s *SampleController) GetFilteredSamples(ctx *gin.Context) {
	filterCheck := ctx.Query("where")
	if samples, found := GetSampleCache(filterCheck); found {
		ctx.JSON(http.StatusOK, gin.H{"data": samples})
		return
	}
	if filterCheck == "" {
		filter := bson.M{}
		samples, err := s.SampleRepository.GetAllFilteredSamples(ctx, filter)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": samples})
		return
	} else {
		whereBsonMap, err := helpers.GetBsonMapFromJsonString(filterCheck)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
		samples, err := s.SampleRepository.GetAllFilteredSamples(ctx, whereBsonMap)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
		SetSampleCache(filterCheck, samples)
		ctx.JSON(http.StatusOK, gin.H{"data": samples})
		return
	}
}

func (s *SampleController) GetOneSampleById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	if sample, found := GetSampleCache(strconv.Itoa(id)); found {
		ctx.JSON(http.StatusOK, gin.H{"data": sample})
		return
	}
	sample, err := s.SampleRepository.GetOneSampleById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
	SetSampleCache(strconv.Itoa(id), sample)
	ctx.JSON(http.StatusOK, gin.H{"data": sample})
}

func (s *SampleController) GetFilteredSampleCount(ctx *gin.Context) {
	jsonString := ctx.Query("where")
	var bsonMap bson.M
	if jsonString == "" {
		bsonMap = bson.M{}
	} else {
		err := json.Unmarshal([]byte(jsonString), &bsonMap)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
	}
	count, err := s.SampleRepository.GetFilteredSampleCount(ctx, bsonMap)
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"data": count})
}

func (s *SampleController) BulkInsertSamples(ctx *gin.Context) {
	var samples []interface{}
	if err := ctx.ShouldBindJSON(&samples); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
	insertedIds, err := s.SampleRepository.InsertMultipleSamples(ctx, samples)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Insertedids": insertedIds})
}

func (s *SampleController) UpdateSingleSamples(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	bsonMap := bson.M{"id": id}
	var update interface{}
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	var updateDoc = make(map[string]interface{})
	updateDoc["$set"] = update
	updateBsonMap, err := helpers.GetBsonMapFromMapDataType(updateDoc)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	updatedCount, err := s.SampleRepository.FindOneSampleAndUpdate(ctx, bsonMap, updateBsonMap)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Updated Count": updatedCount})
}

func (s *SampleController) UpdateMultipleSamples(ctx *gin.Context) {
	filterCheck := ctx.Query("where")
	var update interface{}
	if err := ctx.ShouldBindJSON(&update); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	var updateDoc = make(map[string]interface{})
	updateDoc["$set"] = update
	updateBsonMap, err1 := helpers.GetBsonMapFromMapDataType(updateDoc)
	if err1 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err1})
	}
	var whereBsonMap bson.M
	if filterCheck == "" {
		ctx.JSON(http.StatusBadRequest, errors.New("required where filter"))
		return
	}
	err := json.Unmarshal([]byte(filterCheck), &whereBsonMap)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	updatedCount, err := s.SampleRepository.FindAndUpdateManySamples(ctx, whereBsonMap, updateBsonMap)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
	}
	ctx.JSON(http.StatusAccepted, gin.H{"updatedCount": updatedCount})
}

func (s *SampleController) DeleteSampleById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	deletedCount, err := s.SampleRepository.DeleteSingleSampleById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"Deleted Count": deletedCount})
}
