package api

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ImageUrl string
}

type ImageReq struct {
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

type GetImage struct {
	ID        uint      `json:"id"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type RouterOptions struct {
	Cloudinary *cloudinary.Cloudinary
	GormDB     *gorm.DB
}

func (r *RouterOptions) UploadFile(ctx *gin.Context) {
	var image ImageReq

	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Error: err.Error(),
		})
		return
	}

	if err = ctx.ShouldBind(&image); err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseError{
			Error: err.Error(),
		})
		return
	}
	imageReader, err := image.Image.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Error: err.Error(),
		})
		return
	}

	// id := uuid.New()
	// fileName := id.String() + filepath.Ext(image.Image.Filename)
	// uploading image to cloudinary.com
	resp, err := r.Cloudinary.Upload.Upload(context.Background(), imageReader, uploader.UploadParams{Folder: "samples"})
	// log.Println(fileName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Error: err.Error(),
		})
		return
	}

	// inserting to database with image url uploaded
	req := Image{
		ImageUrl: resp.SecureURL,
	}
	res := r.GormDB.Create(&req)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, GetImage{
		ID:        req.ID,
		ImageUrl:  req.ImageUrl,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		DeletedAt: req.DeletedAt.Time,
	})
}

func (r *RouterOptions) GetFile(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ResponseError{
			Error: err.Error(),
		})
		return
	}
	var image Image
	result := r.GormDB.Find(&image, id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, ResponseError{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, GetImage{
		ID:        image.ID,
		ImageUrl:  image.ImageUrl,
		CreatedAt: image.CreatedAt,
		UpdatedAt: image.UpdatedAt,
		DeletedAt: image.DeletedAt.Time,
	})
}

func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	router.Static("/images", "https://res.cloudinary.com/dmgdx6d2k/image/upload/v1673498952/samples/")

	router.POST("/upload-file", opt.UploadFile)
	router.GET("/file/:id", opt.GetFile)

	return router
}
