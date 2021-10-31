package config

func JaegerCollectorEndpoint() string {
	return GetString("jaeger.reporter.collectorEndpoint")
}

func JaegerServiceName() string {
	return GetString("jaeger.servicename")
}
