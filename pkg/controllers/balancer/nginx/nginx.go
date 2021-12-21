package nginx

import (
	"fmt"
	balancerv1alpha1 "github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1"
	"strings"
)

// Server serves for a typical port with a specific reverse proxy.
type Server struct {
	// TODO: the name segment is not necessary, remove it
	name     string
	protocol string
	port     int32
	upstream string // Server.upstream will be processed exactly by an upstream
}

// conf returns the config segment for the key `server` in nginx.conf.
// Example:
// server {
//     listen 80 tcp;
//     proxy_pass upstream_http;
// }
func (s *Server) conf() string {
	var protocol string
	if s.protocol == "udp" {
		protocol = "udp"
	}
	return fmt.Sprintf(`
server {
    listen %d %s;
    proxy_pass %s;
}
`, s.port, protocol, s.upstream)
}

// backend is the endpoint and its weight for load balancing.
type backend struct {
	name   string
	weight int32
}

// upstream acts as the value of the key `proxy_pass` in nginx.conf.
type upstream struct {
	name     string
	backends []backend
	port     int32
}

// conf returns the config segment for the key `upstream` in nginx.conf.
// Example：
// upstream upstream_http {
//     server example-balancer-v1-backend:80 weight=20;
//     server example-balancer-v2-backend:80 weight=80;
// }
func (us *upstream) conf() string {
	backendStr := ""
	for _, b := range us.backends {
		backendStr += fmt.Sprintf("    server %s:%d weight=%d;\n", b.name, us.port, b.weight)
	}
	return fmt.Sprintf(`
upstream %s {
%s
}
`, us.name, backendStr)
}

// NewConfig generates the `nginx.conf` with the given Balancer instance.
// Example:
// ===================== nginx.conf =====================
// events {
//     worker_connections 1024;
// }
// stream {
//     server {
//         listen 80 tcp;
//         proxy_pass upstream_http;
//     }
//     upstream upstream_http {
//         server example-balancer-v1-backend:80 weight=20;
//         server example-balancer-v2-backend:80 weight=80;
//     }
// }
// ======================================================
func NewConfig(balancer *balancerv1alpha1.Balancer) string {
	var servers []Server
	// TODO: there should be only one element in balancer.Spec.Ports. Check it!
	//  Maybe we could add `+kubebuilder:validation:MaxItems=1` to `Ports` in `BalancerSpec`?
	for _, balancerPort := range balancer.Spec.Ports {
		servers = append(servers, Server{
			name:     balancerPort.Name,
			protocol: strings.ToLower(string(balancerPort.Protocol)),
			port:     int32(balancerPort.Port),
			upstream: fmt.Sprintf("upstream_%s", balancerPort.Name),
		})
	}

	var backends []backend
	for _, balancerBackend := range balancer.Spec.Backends {
		backends = append(backends, backend{
			name:   fmt.Sprintf("%s-%s-backend", balancer.Name, balancerBackend.Name),
			weight: balancerBackend.Weight,
		})
	}

	var upstreams []upstream
	for _, s := range servers {
		upstreams = append(upstreams, upstream{
			name:     s.upstream,
			backends: backends,
			port:     s.port,
		})
	}

	conf := ""
	conf += "events {\n"
	conf += "    worker_connections 1024;\n"
	conf += "}\n"
	conf += "stream {\n"

	for _, s := range servers {
		conf += s.conf()
	}

	for _, us := range upstreams {
		conf += us.conf()
	}

	conf += "}\n"

	return conf
}
