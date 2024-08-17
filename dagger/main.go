// A generated module for FernReporter functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guidewire/fern-reporter/dagger/internal/dagger"
)

type FernReporter struct{}

var platforms = []dagger.Platform{
	"linux/amd64",
	"linux/arm64",
}

// Publish the application container after building and testing it on-the-fly
func (m *FernReporter) Publish(ctx context.Context,
	source *dagger.Directory,
	registry string,
	username string,
	password *dagger.Secret,
) (string, error) {
	_, err := m.Test(ctx, source)
	if err != nil {
		return "", err
	}
	return m.Build(source).
		WithRegistryAuth(registry, username, password).
		// WithExec([]string{"docker", "image", "ls"}).
		Publish(ctx, fmt.Sprintf("%s/%s/my-dagger-test:latest", registry, username)) //#nosec
}

// Build the application container
func (m *FernReporter) Build(source *dagger.Directory) *dagger.Container {
	build := m.BuildEnv(source).
		WithWorkdir("/src").
		WithEnvVariable("CGO_ENABLED", "0").
		// WithExec([]string{"go", "install", "github.com/goreleaser/goreleaser/v2@latest"}).
		WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "fern-reporter", "."}).
		Directory("/src")
	return dag.Container(dagger.ContainerOpts{Platform: "linux/arm64"}).
		From("golang:1.22-alpine").
		WithDirectory("./", build).
		WithExec([]string{"apk", "--no-cache", "add", "ca-certificates"}).
		WithExec([]string{"update-ca-certificates"}).
		WithExec([]string{"ls", "-ltra", "./"})
}

// Return the result of running unit tests
func (m *FernReporter) Test(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.BuildEnv(source).
		WithWorkdir("/src").
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithExec([]string{"go", "install", "github.com/onsi/ginkgo/v2/ginkgo"}).
		WithExec([]string{"ginkgo", "-r"}).
		Stdout(ctx)
}

// Build a ready-to-use development environment
func (m *FernReporter) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1.22-alpine").
		WithMountedDirectory("/src", source)
}

func (m *FernReporter) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}
