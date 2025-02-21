module github.com/kv1sidisi/shapeup/services/usrdatasvc

go 1.23.4

require github.com/kv1sidisi/shapeup/pkg/logger v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/logger => ../../pkg/logger

require (
	github.com/kv1sidisi/shapeup/pkg/config v0.0.0
	github.com/kv1sidisi/shapeup/pkg/database/pgcl v0.0.0
)

replace github.com/kv1sidisi/shapeup/pkg/config => ../../pkg/config

replace github.com/kv1sidisi/shapeup/pkg/utils => ../../pkg/utils

replace github.com/kv1sidisi/shapeup/pkg/database/pgcl => ../../pkg/database/pgcl

require (
	github.com/jackc/pgx/v4 v4.18.3
	github.com/kv1sidisi/shapeup/pkg/errdefs v0.0.0
	google.golang.org/grpc v1.70.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kv1sidisi/shapeup/pkg/utils v0.0.0
	golang.org/x/crypto v0.30.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/kv1sidisi/shapeup/pkg/errdefs => ../../pkg/errdefs

require github.com/kv1sidisi/shapeup/pkg/interceptors/authincp v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/interceptors/authincp => ../../pkg/interceptors/authincp

require (
	github.com/google/uuid v1.6.0
	github.com/kv1sidisi/shapeup/pkg/proto v0.0.0
)

replace github.com/kv1sidisi/shapeup/pkg/proto => ../../pkg/proto
