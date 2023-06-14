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
	Title  string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type createBlogResponse struct {
	BlogID uuid.UUID `json:"blog_id"`
	PairID int64 `json:"pair_id"`
	Title string `json:"title"`
	Content string `json:"content"`
	ImageUrl string `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

func newBlogResponse(blog db.Blog, imageUrl string) createBlogResponse {
	return createBlogResponse{
		BlogID: blog.ID,
		PairID: blog.PairID,
		Title: blog.Title,
		Content: blog.Content,
		ImageUrl: imageUrl,
		CreatedAt: blog.CreateTime,
	}
}


// @Summary      Create Blog
// @Description  Create a new blog
// @Tags         blogs
// @Accept mpfd
// @Param        blog_info body api.createBlogRequest true "Create Blog"
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
		ID: BlogID,
		UserID: authPayload.UserID,
		PairID: User.PairID.Int64,
		Title:  req.Title,
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
