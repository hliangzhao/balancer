package nginx

import (
	`fmt`
	balancerv1alpha1 `github.com/hliangzhao/balancer/pkg/apis/balancer/v1alpha1`
	`strings`
)

// Server serves for a typical port with a specific reverse proxy.
type Server struct {
	// TODO: the name segment is not necessary, remove it
	name     string
	protocol string
	port     int32
	// TODO: define IP:Port as a new type `Address`
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
// When serving for multiple balancer ports, the conf would be:
// ===================== nginx.conf =====================
// events {
//     worker_connections 1024;
// }
// stream {
//     server {
//         listen 80 tcp;
//         proxy_pass upstream_http;
//     }
//     server {
//         listen 8080 udp;
//         proxy_pass upstream_udp;
//     }
//     upstream upstream_http {
//         server example-balancer-v1-backend:80 weight=20;
//         server example-balancer-v2-backend:80 weight=80;
//         server example-balancer-udp-backend:80 weight=80;
//     }
//     upstream upstream_udp {
//         server example-balancer-v1-backend:8080 weight=20;
//         server example-balancer-v2-backend:8080 weight=80;
//         server example-balancer-udp-backend:8080 weight=80;
//     }
// }
// ======================================================
//
// However, what we want:
// ===================== nginx.conf =====================
// events {
//     worker_connections 1024;
// }
// stream {
//     server {
//         listen 80 tcp;
//         proxy_pass upstream_http;
//     }
//     server {
//         listen 8080 udp;
//         proxy_pass upstream_udp;
//     }
//     upstream upstream_http {
//         server example-balancer-v1-backend:80 weight=20;
//         server example-balancer-v2-backend:80 weight=80;
//     }
//     upstream upstream_udp {
//         server example-balancer-udp-backend:8080 weight=80;
//     }
// }
// ======================================================
// Check that whether this is correct and fix it.
func NewConfig(balancer *balancerv1alpha1.Balancer) string {
	var servers []Server
	for _, balancerPort := range balancer.Spec.Ports {
		servers = append(servers, Server{
			name:     balancerPort.Name,
			protocol: strings.ToLower(string(balancerPort.Protocol)),
			port:     int32(balancerPort.Port),
			upstream: fmt.Sprintf("upstream_%s", balancerPort.Name),
		})
	}

	var backends []backend
	// TODO: backends may serve for different balancer ports, but here we get all of them for every server!
	for _, balancerBackend := range balancer.Spec.Backends {
		backends = append(backends, backend{
			name:   fmt.Sprintf("%s-%s-backend", balancer.Name, balancerBackend.Name),
			weight: balancerBackend.Weight,
		})
	}

	var upstreams []upstream
	for _, s := range servers {
		// TODO: find backendsForThisServer from backends
		upstreams = append(upstreams, upstream{
			name:     s.upstream,
			backends: backends, // TODO: replaced by backendsForThisServer
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
