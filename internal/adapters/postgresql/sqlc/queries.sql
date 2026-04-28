-- name: ListProducts :many
Select * from products;

-- name: GetProductByID :one
Select * from products where id = $1;