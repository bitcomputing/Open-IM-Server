all: format gen

format: 
	goctl api format --dir ./api

gen:
	goctl api go -api ./api/gateway.api -dir .