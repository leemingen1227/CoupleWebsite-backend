package db

import (
	"context"
)

type CreateInvitationTxParams struct {
	CreateInvitationParams
	AfterCreate func(invitation Invitation) error
}

type CreateInvitationTxResult struct {
	Invitation Invitation
}

func (store *SQLStore) CreateInvitationTx(ctx context.Context, arg CreateInvitationTxParams) (CreateInvitationTxResult, error) {
	var result CreateInvitationTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Invitation, err = q.CreateInvitation(ctx, arg.CreateInvitationParams)

		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Invitation)
	})

	return result, err
}