run:
	go run test.go

get:
	curl -H "Content-Type: application/json" \
	-H "x-api-key: thisisanapikey" \
	http://localhost:3000/containers

get_test:
	curl -H "Content-Type: application/json" \
	http://localhost:3000/test