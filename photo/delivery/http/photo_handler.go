package delivery

import (
	"final-project-mygram-fachmi/domain"
	"final-project-mygram-fachmi/helpers"
	"final-project-mygram-fachmi/photo/delivery/http/middleware"
	"final-project-mygram-fachmi/photo/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type photoHandler struct {
	photoUseCase domain.PhotoUseCase
}

func NewPhotoHandler(routers *gin.Engine, photoUseCase domain.PhotoUseCase) {
	handler := &photoHandler{photoUseCase}

	router := routers.Group("/photos")
	{
		router.Use(middleware.Authentication())
		router.GET("", handler.Fetch)
		router.POST("", handler.Store)
		router.PUT("/:photoId", middleware.Authorization(handler.photoUseCase), handler.Update)
		router.DELETE("/:photoId", middleware.Authorization(handler.photoUseCase), handler.Delete)
	}
}

func (handler *photoHandler) Fetch(ctx *gin.Context) {
	var (
		photos []domain.Photo
		err    error
	)

	if err = handler.photoUseCase.Fetch(ctx.Request.Context(), &photos); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	fetchedPhotos := []*utils.FetchedPhoto{}

	for _, photo := range photos {
		fetchedPhotos = append(fetchedPhotos, &utils.FetchedPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User: &utils.User{
				Email:    photo.User.Email,
				Username: photo.User.Username,
			},
		})
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data:   fetchedPhotos,
	})
}

func (handler *photoHandler) Store(ctx *gin.Context) {
	var (
		photo domain.Photo
		err   error
	)

	userData := ctx.MustGet("userData").(jwt.MapClaims)
	userID := string(userData["id"].(string))

	if err = ctx.ShouldBindJSON(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	photo.UserID = userID

	if err = handler.photoUseCase.Store(ctx.Request.Context(), &photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, helpers.ResponseData{
		Status: "success",
		Data: utils.AddedPhoto{
			ID:        photo.ID,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			UserID:    photo.UserID,
			CreatedAt: photo.CreatedAt,
		},
	})
}

func (handler *photoHandler) Update(ctx *gin.Context) {
	var (
		photo domain.Photo
		err   error
	)

	if err = ctx.ShouldBindJSON(&photo); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	updatedPhoto := domain.Photo{
		Title:    photo.Title,
		Caption:  photo.Caption,
		PhotoUrl: photo.PhotoUrl,
	}

	photoID := ctx.Param("photoId")

	if photo, err = handler.photoUseCase.Update(ctx.Request.Context(), updatedPhoto, photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, helpers.ResponseData{
		Status: "success",
		Data: utils.UpdatedPhoto{
			ID:        photo.ID,
			UserID:    photo.UserID,
			Title:     photo.Title,
			PhotoUrl:  photo.PhotoUrl,
			Caption:   photo.Caption,
			UpdatedAt: photo.UpdatedAt,
		},
	})
}

func (handler *photoHandler) Delete(ctx *gin.Context) {
	photoID := ctx.Param("photoId")

	if err := handler.photoUseCase.Delete(ctx.Request.Context(), photoID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helpers.ResponseMessage{
			Status:  "fail",
			Message: err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "your photo has been successfully deleted",
	})
}
