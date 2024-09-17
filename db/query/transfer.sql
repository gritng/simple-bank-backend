-- name: CreateTransfer :one
INSERT INTO transfers (
  from_acc_id, 
  to_acc_id, 
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfer :many
SELECT * FROM transfers
ORDER BY id
LIMIT $1 
OFFSET $2;

-- name: UpdateTransfer :exec
UPDATE transfers
  set from_acc_id = $2,
    to_acc_id = $3,
    amount = $4
WHERE id = $1;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;