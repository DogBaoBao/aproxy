package pkg

// HealthCheck
type HealthCheck struct {
}

// HttpHealthCheck
type HttpHealthCheck struct {
	HealthCheck
	Host             string
	Path             string
	UseHttp2         bool
	ExpectedStatuses int64
}

type GrpcHealthCheck struct {
	HealthCheck
	ServiceName string
	Authority   string
}

type CustomHealthCheck struct {
	HealthCheck
	Name   string
	Config interface{}
}
