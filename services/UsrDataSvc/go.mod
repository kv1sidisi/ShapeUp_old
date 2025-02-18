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

require github.com/kv1sidisi/shapeup/pkg/errdefs v0.0.0

replace github.com/kv1sidisi/shapeup/pkg/errdefs => ../../pkg/errdefs