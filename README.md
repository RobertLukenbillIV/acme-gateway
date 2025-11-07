# acme-gateway
acme-gateway is the entry point to the Acme platform. Built in Go for performance and simplicity, it authenticates and authorizes incoming requests by verifying JWTs from acme-auth-service, then proxies them to internal services such as acme-tickets-service.
