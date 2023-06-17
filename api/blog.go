package api

import (
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leemingen1227/couple-server/aws"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"github.com/leemingen1227/couple-server/token"
)

type createBlogRequest struct {
	Title   string                `form:"title" binding:"required"`
	Content string                `form:"content" binding:"required"`
	Image   *multipart.FileHeader `form:"image"`
}

type createBlogResponse struct {
	BlogID    uuid.UUID `json:"blog_id"`
	PairID    int64     `json:"pair_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

func newBlogResponse(blog db.Blog, imageUrl string) createBlogResponse {
	return createBlogResponse{
		BlogID:    blog.ID,
		PairID:    blog.PairID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageUrl:  imageUrl,
		CreatedAt: blog.CreateTime,
	}
}

// @Summary      Create Blog
// @Description  Create a new blog
// @Tags         blogs
// @Param        Authorization     header    string     true   "Bearer token"
// @Param title formData string true "Blog Title"
// @Param content formData string true "Blog Content"
// @Param image formData file true "Blog Image"
// @Success      200  {object}  api.createBlogResponse
// @Router       /blogs [post]
func (server *Server) createBlog(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Printf("req: %v", req)

	//Get the user information from the context
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	sess, err := aws.InitAWS()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var imageId string
	if req.Image != nil {
		imageId, err = aws.UploadImageToS3(sess, req.Image)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	BlogID, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	User, err := server.store.GetUser(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	blog, err := server.store.CreateBlog(ctx, db.CreateBlogParams{
		ID:      BlogID,
		UserID:  authPayload.UserID,
		PairID:  User.PairID.Int64,
		Title:   req.Title,
		Content: req.Content,
		Picture: imageId,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Assuming a function server.s3.GetSignedURL() that generates a signed URL for an S3 object.
	var imageUrl string
	if imageId != "" {
		imageUrl, err = aws.GetSignedURL(sess, imageId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	rsp := newBlogResponse(blog, imageUrl)
	ctx.JSON(http.StatusOK, rsp)
}

type getBlogByBlogIDRequest struct {
	blogID uuid.UUID `uri:"blogID" binding:"required"`
}

type getBlogByBlogIDResponse struct {
	BlogID    uuid.UUID `json:"blog_id"`
	PairID    int64     `json:"pair_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

func newGetBlogByBlogIDResponse(blog db.Blog, imageUrl string) getBlogByBlogIDResponse {
	return getBlogByBlogIDResponse{
		BlogID:    blog.ID,
		PairID:    blog.PairID,
		Title:     blog.Title,
		Content:   blog.Content,
		ImageUrl:  imageUrl,
		CreatedAt: blog.CreateTime,
	}
}

// @Summary      Get Blog
// @Description  Get a blog
// @Tags         blogs
// @Param        Authorization     header    string     true   "Bearer token"
// @Param blogID path string true "Blog ID"
// @Success      200  {object}  api.getBlogByBlogIDResponse
// @Router       /blogs/blog/{blogID} [get]
func (server *Server) getBlogByBlogID(ctx *gin.Context) {
	blogIDStr := ctx.Param("blogID")
	if blogIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Blog ID required"})
		return
	}

	//parse the blog id into uuid, since when binding from URI or JSON, Gin treats UUIDs as strings.
	blogID, err := uuid.Parse(blogIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Blog ID format"})
		return
	}

	blog, err := server.store.GetBlogByBlogID(ctx, blogID)
	if err != nil {
		log.Printf("can't get blog: %v", err)
		// log.Printf("blogID: %v", req.blogID)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sess, err := aws.InitAWS()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var imageUrl string
	if blog.Picture != "" {
		imageUrl, err = aws.GetSignedURL(sess, blog.Picture)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	rsp := newGetBlogByBlogIDResponse(blog, imageUrl)
	ctx.JSON(http.StatusOK, rsp)
}

type getBlogsByPairIDRequest struct {
	PairID   int64 `uri:"pairID" binding:"required"`
	Page     int   `form:"page" binding:"omitempty"`
	PageSize int   `form:"page_size" binding:"omitempty"`
}

type getBlogsByPairIDResponse struct {
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
	Blogs    []getBlogByBlogIDResponse `json:"blogs"`
}

// @Summary      Get Blogs by PairID
// @Description  Get blogs by pair id
// @Tags         blogs
// @Param        Authorization     header    string     true   "Bearer token"
// @Param        pairID           path      int        true   "Pair ID"
// @Param        page              query     int        false   "Page number"
// @Param        page_size         query     int        false   "Page size"
// @Success      200  {object}  api.getBlogsByPairIDResponse
// @Router       /blogs/{pairID} [get]
func (server *Server) getBlogsByPairID(ctx *gin.Context) {
	var req getBlogsByPairIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil { // Bind Page and PageSize from query parameters
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// set default values if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 3
	}

	total, err := server.store.CountBlogsByPairID(ctx, req.PairID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	offset := (req.Page - 1) * req.PageSize
	blogs, err := server.store.GetBlogsByPairID(ctx, db.GetBlogsByPairIDParams{
		PairID: req.PairID,
		Limit:  int32(req.PageSize),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	sess, err := aws.InitAWS()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	blogs_responses := make([]getBlogByBlogIDResponse, len(blogs))
	for i, blog := range blogs {
		imageUrl, err := aws.GetSignedURL(sess, blog.Picture)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		blogs_responses[i] = newGetBlogByBlogIDResponse(blog, imageUrl)
	}

	rsp := getBlogsByPairIDResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Blogs:    blogs_responses,
	}
	ctx.JSON(http.StatusOK, rsp)
}
