module github.com/kv1sidisi/shapeup/pkg/interceptors/authincp

go 1.23.4

require (
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
)

require github.com/kv1sidisi/shapeup/pkg/errdefs v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/errdefs => ../../errdefs

require github.com/kv1sidisi/shapeup/pkg/proto v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/proto => ../../proto
