package controller_manager

import (
	`sigs.k8s.io/controller-runtime/pkg/log`
)

var (
	metricsHost         = "0.0.0.0"
	metricsPort         = 8383
	operatorMetricsPort = 8686
)

var logger = log.Log.WithName("cmd")

