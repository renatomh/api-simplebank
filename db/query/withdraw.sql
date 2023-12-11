-- name: CreateWithdraw :one
INSERT INTO withdraws (
  account_id,
  amount,
  "user"
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: Getwithdraw :one
SELECT * FROM withdraws
WHERE id = $1 LIMIT 1;

-- name: ListWithdraws :many
SELECT * FROM withdraws
WHERE 
    account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
