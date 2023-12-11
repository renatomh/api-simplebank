-- name: CreateDeposit :one
INSERT INTO deposits (
  account_id,
  amount,
  "user"
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetDeposit :one
SELECT * FROM deposits
WHERE id = $1 LIMIT 1;

-- name: ListDeposits :many
SELECT * FROM deposits
WHERE 
    account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;
