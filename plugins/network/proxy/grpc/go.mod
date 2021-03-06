module github.com/asim/go-micro/plugins/network/proxy/grpc/v3

go 1.15

require (
	github.com/asim/go-micro/plugins/client/grpc/v3 v3.0.0-00010101000000-000000000000
	github.com/asim/go-micro/v3 v3.0.0-20210120135431-d94936f6c97c
)

replace github.com/asim/go-micro/plugins/registry/memory/v3 => ../../../registry/memory

replace github.com/asim/go-micro/plugins/client/grpc/v3 => ../../../client/grpc
