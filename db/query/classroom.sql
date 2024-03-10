-- name: GetClassrooms: many
select *
from classroom;

-- name: GetClassroom: one
select *
from classroom
where id = $1;

