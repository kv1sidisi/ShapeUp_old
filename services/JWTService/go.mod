module github.com/kv1sidisi/shapeup/services/jwtsvc

go 1.23.4

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/ilyakaznacheev/cleanenv v1.5.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

require github.com/kv1sidisi/shapeup/pkg/logger v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/logger => ../../pkg/logger

require github.com/kv1sidisi/shapeup/pkg/config v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/config => ../../pkg/config

require github.com/kv1sidisi/shapeup/pkg/errdefs v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/errdefs => ../../pkg/errdefs

require (
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/kv1sidisi/shapeup/pkg/proto v0.0.0
	google.golang.org/grpc v1.70.0
)

replace github.com/kv1sidisi/shapeup/pkg/proto => ../../pkg/proto
