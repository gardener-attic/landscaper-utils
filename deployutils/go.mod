module github.com/gardener/landscaper-utils/deployutils

go 1.16

require (
	github.com/gardener/landscaper/apis v0.15.1
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/spf13/pflag v1.0.5
	go.uber.org/zap v1.19.1
	k8s.io/client-go v0.22.3
	sigs.k8s.io/controller-runtime v0.10.2
	sigs.k8s.io/yaml v1.3.0
)
