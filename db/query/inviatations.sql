-- name: CreateInvitation :one
INSERT INTO invitations (inviter_id, invitee_email, invitation_token)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateInvitation :one
UPDATE invitations
SET is_accepted = true
WHERE id = $1
AND invitation_token = $2
AND is_accepted = false
RETURNING *;

-- name: GetInvitation :one
SELECT * FROM invitations
WHERE id = $1 LIMIT 1;