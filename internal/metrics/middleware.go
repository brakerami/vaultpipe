package metrics

// SecretFetched records the outcome of a Vault secret fetch operation
// into the provided Registry using the well-known counter names defined
// in the package documentation.
func SecretFetched(r *Registry, err error) {
	if err != nil {
		r.Counter("secret.fetch.error").Inc()
		return
	}
	r.Counter("secret.fetch.ok").Inc()
}

// CacheAccess records a cache hit or miss into the provided Registry.
func CacheAccess(r *Registry, hit bool) {
	if hit {
		r.Counter("cache.hit").Inc()
		return
	}
	r.Counter("cache.miss").Inc()
}

// LeaseRenewed records the outcome of a lease renewal attempt.
func LeaseRenewed(r *Registry, err error) {
	if err != nil {
		r.Counter("renew.error").Inc()
		return
	}
	r.Counter("renew.ok").Inc()
}
