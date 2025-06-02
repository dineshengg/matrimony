module webserver

go 1.24.3

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/lib/pq v1.10.9
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dineshengg/matrimony v0.0.0-00010101000000-000000000000 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.37.0 // indirect
	github.com/qiangxue/fasthttp-routing v0.0.0-20160225050629-6ccdc2a18d87 // indirect
	github.com/redis/go-redis/v9 v9.7.3 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.62.0 // indirect
)

replace github.com/dineshengg/matrimony => ./../../Matrimony
replace github.com/dineshengg/matrimony/common => ./../../Matrimony/common
replace github.com/dineshengg/matrimony/middleware => ./../../Matrimony/middleware
replace github.com/dineshengg/matrimony/userprofile/login => ./../../Matrimony/UserProfile/login
