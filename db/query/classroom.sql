-- name: ListClassrooms :many
select room_code, name, description, category_id, image_urls
from classroom;

-- name: GetClassroom :one
select *
from classroom
where id = $1;

