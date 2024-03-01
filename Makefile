
run:
	go run main.go


gen:
	sqlc generate

shuv:
	git add .
	git commit -a
	git push

.PHONY: run, gen, shuv
