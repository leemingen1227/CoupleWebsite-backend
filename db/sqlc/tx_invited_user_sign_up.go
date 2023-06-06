package db

import (
	"context"
	"database/sql"
	"log"
	"time"
	"github.com/google/uuid"
)

type InvitedUserSignUpTxParams struct {
	InvitationID int64
	InvitationToken string
	CreateUserParams
	AfterCreate func(user User) error
}

type InvitedUserSignUpTxResult struct {
	User User
	Invitation Invitation
}

func (store *SQLStore) InvitedUserSignUpTx(ctx context.Context, arg InvitedUserSignUpTxParams) (InvitedUserSignUpTxResult, error) {
	var result InvitedUserSignUpTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Invitation, err = q.UpdateInvitation(ctx, UpdateInvitationParams{
			ID: arg.InvitationID,
			InvitationToken : arg.InvitationToken,
		})
		if err != nil {
			log.Print("Error in UpdateInvitation")
			return err
		}
		//create new pair record
		var pair Pair
		pair, err = q.CreatePair(ctx, sql.NullTime{
			Time: time.Time{},
			Valid: false,
		})
		if err != nil {
			log.Print("Error in CreatePair")
			return err
		}

		//update inviter pair id in user table
		var inviter User
		inviter, err = q.UpdateUser(ctx, UpdateUserParams{
			ID: result.Invitation.InviterID,
			PairID: sql.NullInt64{
				Int64: pair.ID,
				Valid: true,
			},
		})
		if err != nil {
			log.Print("Error in UpdateInviter")
			return err
		}

		var uid uuid.UUID
		uid, err = uuid.NewRandom()
		if err != nil {
			return err
		}

		//add pair id to CreateUserParams
		arg.CreateUserParams.PairID = sql.NullInt64{
			Int64: pair.ID,
			Valid: true,
		}

		//add uuid to CreateUserParams
		arg.CreateUserParams.ID = uid
		
		// update the CreateUserParams isEmailVerified field to true
		

		// Create the invited user with the new pair id
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			log.Print("Error in CreateUser")
			return err
		}

		//create two UserPair records for the new pair
		_, err = q.CreateUsersPair(ctx, CreateUsersPairParams{
			UserID: result.User.ID,
			PairID: pair.ID,
		})
		if err != nil {
			log.Print("Error in CreateUserPair")
			return err
		}
		_, err = q.CreateUsersPair(ctx, CreateUsersPairParams{
			UserID: inviter.ID,
			PairID: pair.ID,
		})
		if err != nil {
			log.Print("Error in CreateUserPair")
			return err
		}

		return err

	})
	return result, err
}