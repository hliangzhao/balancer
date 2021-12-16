package main

import (
	`context`
	`flag`
	`fmt`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	`github.com/hliangzhao/balancer/pkg/controllers`
	`github.com/operator-framework/operator-lib/leader`
	`os`
	`sigs.k8s.io/controller-runtime/pkg/client/config`
	logf `sigs.k8s.io/controller-runtime/pkg/log`
	`sigs.k8s.io/controller-runtime/pkg/log/zap`
	`sigs.k8s.io/controller-runtime/pkg/manager`
)

const (
	metricsHost         = "0.0.0.0"
	metricsPort         = 8383
	operatorMetricsPort = 8686
)

var log = logf.Log.WithName("cmd")

func main() {
	// define the custom zap logger
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	logger := zap.New(zap.UseFlagOptions(&opts))
	logf.SetLogger(logger)

	// get a config to talk to the api-server
	log.Info("Getting kubeconfig.")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Failed to get kubeconfig")
		os.Exit(1)
	}

	ctx := context.Background()

	// become the leader before proceeding
	err = leader.Become(ctx, "balancer-lock")
	if err != nil {
		log.Error(err, "Error happened in leader election")
		os.Exit(1)
	}

	// setup controller-manager with metrics serving
	log.Info("Setting up controller-manager.")
	watchNs, err := GetWatchNamespace()
	if err != nil {
		log.Error(err, "Failed to get watch namespace")
		os.Exit(1)
	}
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          watchNs,
		MapperProvider:     NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		log.Error(err, "Failed to set up manager")
		os.Exit(1)
	}

	// setup Scheme for all resources
	log.Info("Registering Components.")
	if err := balancerv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "Failed to setup scheme for all resources")
		os.Exit(1)
	}

	// setup all Controllers
	log.Info("Setting up all controllers.")
	if err := controllers.AddToManager(mgr); err != nil {
		log.Error(err, "Failed to set up controllers")
		os.Exit(1)
	}

	// TODO: serve custom metrics
}