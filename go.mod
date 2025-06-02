module github.com/dineshengg/matrimony

go 1.24.3

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/lib/pq v1.10.9
	github.com/qiangxue/fasthttp-routing v0.0.0-20160225050629-6ccdc2a18d87
	github.com/redis/go-redis/v9 v9.9.0
	github.com/valyala/fasthttp v1.62.0
	golang.org/x/crypto v0.38.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-ozzo/ozzo-routing v2.1.4+incompatible // indirect
	github.com/golang/gddo v0.0.0-20210115222349-20d68f94ee1f // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace github.com/dineshengg/matrimony => ./../Matrimony

replace github.com/dineshengg/matrimony/common => ./../Matrimony/common

replace github.com/dineshengg/matrimony/common/utils => ./../Matrimony/common/utils
