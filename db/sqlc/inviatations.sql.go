// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0
// source: inviatations.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createInvitation = `-- name: CreateInvitation :one
INSERT INTO invitations (inviter_id, invitee_email, invitation_token)
VALUES ($1, $2, $3)
RETURNING id, inviter_id, invitee_email, invitation_token, is_accepted, create_time
`

type CreateInvitationParams struct {
	InviterID       uuid.UUID `json:"inviter_id"`
	InviteeEmail    string    `json:"invitee_email"`
	InvitationToken string    `json:"invitation_token"`
}

func (q *Queries) CreateInvitation(ctx context.Context, arg CreateInvitationParams) (Invitation, error) {
	row := q.db.QueryRowContext(ctx, createInvitation, arg.InviterID, arg.InviteeEmail, arg.InvitationToken)
	var i Invitation
	err := row.Scan(
		&i.ID,
		&i.InviterID,
		&i.InviteeEmail,
		&i.InvitationToken,
		&i.IsAccepted,
		&i.CreateTime,
	)
	return i, err
}

const getInvitation = `-- name: GetInvitation :one
SELECT id, inviter_id, invitee_email, invitation_token, is_accepted, create_time FROM invitations
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetInvitation(ctx context.Context, id int64) (Invitation, error) {
	row := q.db.QueryRowContext(ctx, getInvitation, id)
	var i Invitation
	err := row.Scan(
		&i.ID,
		&i.InviterID,
		&i.InviteeEmail,
		&i.InvitationToken,
		&i.IsAccepted,
		&i.CreateTime,
	)
	return i, err
}

const updateInvitation = `-- name: UpdateInvitation :one
UPDATE invitations
SET is_accepted = true
WHERE id = $1
AND invitation_token = $2
AND is_accepted = false
RETURNING id, inviter_id, invitee_email, invitation_token, is_accepted, create_time
`

type UpdateInvitationParams struct {
	ID              int64  `json:"id"`
	InvitationToken string `json:"invitation_token"`
}

func (q *Queries) UpdateInvitation(ctx context.Context, arg UpdateInvitationParams) (Invitation, error) {
	row := q.db.QueryRowContext(ctx, updateInvitation, arg.ID, arg.InvitationToken)
	var i Invitation
	err := row.Scan(
		&i.ID,
		&i.InviterID,
		&i.InviteeEmail,
		&i.InvitationToken,
		&i.IsAccepted,
		&i.CreateTime,
	)
	return i, err
}
