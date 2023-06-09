// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package godata loads the Go project's corpus of Git, Github, and
// Gerrit activity into memory to allow easy analysis without worrying
// about APIs and their pagination, quotas, and other nuisances and
// limitations.
package godata

import (
	"context"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"golang.org/x/build/maintner"
)

// Server is the Go project's production maintner log.
const Server = "https://maintner.golang.org/logs"

// Get returns the Go project's corpus, containing all Git commits,
// Github activity, and Gerrit activity and metadata since the
// beginning of the project.
//
// Use Corpus.Update to keep the corpus up-to-date. If you do this, you must
// hold the read lock if reading and updating concurrently.
//
// The initial call to Get will download a few gigabytes of data
// into a directory "golang-maintner" under your operating
// system's user cache directory. Subsequent calls will only download
// what's changed since the previous call.
//
// Even with all the data already cached on local disk, a call to Get
// takes approximately 15 seconds per gigabyte of mutation log data
// to load it into memory.
// For daemons, use Corpus.Update to incrementally update an
// already-loaded Corpus.
//
// The in-memory representation is about 25% larger than its on-disk
// size. In April 2022, it's under 4 GB.
//
// See https://pkg.go.dev/golang.org/x/build/maintner#Corpus for how
// to walk the data structure.
func Get(ctx context.Context) (*maintner.Corpus, error) {
	targetDir := Dir()
	if err := os.MkdirAll(targetDir, 0700); err != nil {
		return nil, err
	}
	mutSrc := maintner.NewNetworkMutationSource(Server, targetDir)
	corpus := new(maintner.Corpus)
	if err := corpus.Initialize(ctx, mutSrc); err != nil {
		return nil, err
	}
	return corpus, nil
}

// Dir returns the directory containing the cached mutation logs.
func Dir() string {
	return filepath.Join(XdgCacheDir(), "golang-maintner")
}

// XdgCacheDir returns the XDG Base Directory Specification cache
// directory.
func XdgCacheDir() string {
	cache := os.Getenv("XDG_CACHE_HOME")
	if cache != "" {
		return cache
	}
	home := homeDir()
	// Not XDG but standard for OS X.
	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library/Caches")
	}
	return filepath.Join(home, ".cache")
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	home := os.Getenv("HOME")
	if home != "" {
		return home
	}
	u, err := user.Current()
	if err != nil {
		log.Fatalf("failed to get home directory or current user: %v", err)
	}
	return u.HomeDir
}
