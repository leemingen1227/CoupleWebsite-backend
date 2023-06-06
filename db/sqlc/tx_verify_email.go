package db

import (
	"context"
	"database/sql"
	"log"
)

type VerifyEmailTxParams struct {
	EmailId int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
    var result VerifyEmailTxResult

    err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID: arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			log.Print("Error in UpdateVerifyEmail")
			return err
		}

		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			ID: result.VerifyEmail.UserID,
			IsEmailVerified: sql.NullBool{
				Bool: true,
				Valid: true,
			},
		})
		if err != nil {
			log.Print("Error in UpdateUser")
			return err
		}
		
		return err
    })

    return result, err
}