# run in container
go build ./cmd/service
nohup ./service > service.log & exit