package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/leemingen1227/couple-server/db/sqlc"
	"net/http"
	"strconv"
	// "log"
	// "io/ioutil"
	// "encoding/json"
	// "fmt"
)

type VerifyEmailRequest struct {
	EmailId int64 `json:"email_id" binding:"required"`
	SecretCode string `json:"secret_code" binding:"required"`
}

type VerifyEmailResponse struct {
	IsEmailVerified bool `json:"is_email_verified"`
}

func newVerifyEmailResponse(verify_email_tx_result db.VerifyEmailTxResult) VerifyEmailResponse {
	return VerifyEmailResponse{
		IsEmailVerified: verify_email_tx_result.User.IsEmailVerified,
	}
}

//the handler for verifying email.
// the handler for verify email.
// @Summary      Verify email
// @Description  Verify the email of created account.
// @Tags         verify
// @Param		 verify_info	query	api.VerifyEmailRequest	true	"Verify Email"
// @Success      200  {object}  api.VerifyEmailResponse
// @Router       /verify/verify_email [get]
func (server *Server) verifyEmail(ctx *gin.Context){
	var req VerifyEmailRequest

	// Use the Query method to get a query parameter
	emailId, isPresent := ctx.GetQuery("email_id")
	if !isPresent {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email_id parameter not found"})
		return
	}

	secretCode, isPresent := ctx.GetQuery("secret_code")
	if !isPresent {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "secret_code parameter not found"})
		return
	}

	// Assuming EmailId is an integer
	req.EmailId, _ = strconv.ParseInt(emailId, 10, 64)
	req.SecretCode = secretCode

	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	log.Printf("Error binding request: %v", err) // Log the error
	// 	ctx.JSON(http.StatusBadRequest, errorResponse(err))
	// 	return
	// }

	// bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	// log.Print(ctx.Request.Body)
	// if err != nil {
	// 	log.Printf("Error reading body: %v", err)
		
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "can't read body"})
	// 	return
	// }

	// if err := json.Unmarshal(bodyBytes, &req); err != nil {
	// 	log.Printf("Error unmarshaling body: %v", err)
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body format"})
	// 	return
	// }

	// // Now we can check req's fields and log if they are not as expected
	// if req.EmailId == 0 {
	// 	log.Println("Missing EmailId")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing EmailId"})
	// 	return
	// }
	// if req.SecretCode == "" {
	// 	log.Println("Missing SecretCode")
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing SecretCode"})
	// 	return
	// }

	verify_email_tx_result, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId: req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newVerifyEmailResponse(verify_email_tx_result)
	ctx.JSON(http.StatusOK, rsp)
}

