// Package pipeline wires together the secret resolution, environment
// construction, and process execution stages of vaultpipe into a single
// cohesive unit.
//
// Typical usage:
//
//	cfg, _ := config.LoadFile("vaultpipe.yaml")
//	p, base, err := pipeline.Build(cfg, pipeline.BuildOptions{
//		Args:   os.Args[1:],
//		Logger: logger,
//		Cache:  cache.New(ttl),
//	})
//	if err != nil { ... }
//	if err := p.Run(ctx, base); err != nil { ... }
//
// The Pipeline itself is intentionally dependency-injected so that tests
// can substitute stub resolvers and runners without network access.
package pipeline
