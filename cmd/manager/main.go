/*
Copyright 2021 hliangzhao.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	exposerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	"github.com/hliangzhao/balancer/pkg/controllers"
	"os"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(exposerv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "24a4cdb8.hliangzhao.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = controllers.AddToManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Balancer")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder
	// if err = (&balancer.ReconcilerBalancer{
	// 	Client: mgr.GetClient(),
	// 	Scheme: mgr.GetScheme(),
	// }).SetupWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create controller", "controller", "Balancer")
	// 	os.Exit(1)
	// }

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// TODO: to-be-deleted
// var log = logf.Log.WithName("balancer-cmd")
//
// func main() {
// 	// define the custom zap logger
// 	opts := zap.Options{}
// 	opts.BindFlags(flag.CommandLine)
// 	flag.Parse()
// 	logger := zap.New(zap.UseFlagOptions(&opts))
// 	logf.SetLogger(logger)
//
// 	// get a config to talk to the api-server
// 	log.Info("Getting kubeconfig.")
// 	cfg, err := config.GetConfig()
// 	if err != nil {
// 		log.Error(err, "Failed to get kubeconfig")
// 		os.Exit(1)
// 	}
//
// 	ctx := context.Background()
//
// 	// become the leader before proceeding
// 	err = leader.Become(ctx, "balancer-lock")
// 	if err != nil {
// 		log.Error(err, "Error happened in leader election")
// 		os.Exit(1)
// 	}
//
// 	// setup controller-manager with metrics serving
// 	log.Info("Setting up controller-manager.")
// 	watchNs, err := app.GetWatchNamespace()
// 	if err != nil {
// 		log.Error(err, "Failed to get watch namespace")
// 		os.Exit(1)
// 	}
// 	mgr, err := manager.New(cfg, manager.Options{
// 		Namespace:      watchNs,
// 		MapperProvider: app.NewDynamicRESTMapper,
// 		// MetricsBindAddress is the TCP address that the controller should bind to for serving prometheus metrics
// 		MetricsBindAddress: fmt.Sprintf("%s:%d", app.MetricsHost, app.MetricsPort),
// 	})
// 	if err != nil {
// 		log.Error(err, "Failed to set up manager")
// 		os.Exit(1)
// 	}
//
// 	// setup scheme for all resources
// 	log.Info("Registering Components.")
// 	if err := balancerv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
// 		log.Error(err, "Failed to setup scheme for all resources")
// 		os.Exit(1)
// 	}
//
// 	// setup all controllers
// 	log.Info("Setting up all controllers.")
// 	if err := controllers.AddToManager(mgr); err != nil {
// 		log.Error(err, "Failed to set up controllers")
// 		os.Exit(1)
// 	}
//
// 	// TODO: how to serve custom metrics for v1.22.2?
//
// 	// finally, we start the manager
// 	log.Info("Starting the cmd.")
// 	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
// 		log.Error(err, "Unable start the manager")
// 		os.Exit(1)
// 	}
// }
