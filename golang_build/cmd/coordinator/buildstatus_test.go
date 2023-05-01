// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.16 && (linux || darwin)
// +build go1.16
// +build linux darwin

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/build/buildenv"
	"golang.org/x/build/dashboard"
	"golang.org/x/build/internal/buildgo"
	"golang.org/x/build/internal/coordinator/pool"
)

// TestParseOutputAndHeader tests header parsing by parseOutputAndHeader.
func TestParseOutputAndHeader(t *testing.T) {
	for _, tc := range []struct {
		name         string
		input        []byte
		wantMetadata string
		wantHeader   string
		wantOut      []byte
	}{
		{
			name: "standard",
			input: []byte(`
XXXBANNERXXX:Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: "",
			wantHeader:   "##### Testing packages.",
			wantOut: []byte(`ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "header only",
			input: []byte(`
XXXBANNERXXX:Testing packages.
`),
			wantMetadata: "",
			wantHeader:   "##### Testing packages.",
			wantOut:      []byte(``),
		},
		{
			name: "header only missing trailing newline",
			input: []byte(`
XXXBANNERXXX:Testing packages.`),
			wantMetadata: "",
			wantHeader:   "##### Testing packages.",
			wantOut:      []byte(``),
		},
		{
			name: "no banner",
			input: []byte(`ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: "",
			wantHeader:   "",
			wantOut: []byte(`ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "no newline",
			input: []byte(`XXXBANNERXXX:Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: "",
			wantHeader:   "",
			wantOut: []byte(`XXXBANNERXXX:Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "wrong banner",
			input: []byte(`
##### Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: "",
			wantHeader:   "",
			wantOut: []byte(`
##### Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "metadata",
			input: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz

XXXBANNERXXX:Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: `##### Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz`,
			wantHeader: "##### Testing packages.",
			wantOut: []byte(`ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "metadata missing separator newline",
			input: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
XXXBANNERXXX:Testing packages.
ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
			wantMetadata: `##### Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz`,
			wantHeader: "##### Testing packages.",
			wantOut: []byte(`ok	archive/tar	0.015s
ok	archive/zip	0.406s
ok	bufio	0.075s
`),
		},
		{
			name: "metadata missing second banner",
			input: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
`),
			wantMetadata: "",
			wantHeader:   "",
			wantOut: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz
`),
		},
		{
			name: "metadata missing body",
			input: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz

XXXBANNERXXX:Testing packages.
`),
			wantMetadata: `##### Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz`,
			wantHeader: "##### Testing packages.",
			wantOut:    []byte(``),
		},
		{
			name: "metadata missing body and newline",
			input: []byte(`
XXXBANNERXXX:Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz

XXXBANNERXXX:Testing packages.`),
			wantMetadata: `##### Test execution environment.
# GOARCH: amd64
# CPU: Intel(R) Xeon(R) W-2135 CPU @ 3.70GHz`,
			wantHeader: "##### Testing packages.",
			wantOut:    []byte(``),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			gotMetadata, gotHeader, gotOut := parseOutputAndHeader(tc.input)
			if gotMetadata != tc.wantMetadata {
				t.Errorf("parseOutputAndBanner(%q) got metadata %q want metadata %q", string(tc.input), gotMetadata, tc.wantMetadata)
			}
			if gotHeader != tc.wantHeader {
				t.Errorf("parseOutputAndBanner(%q) got header %q want header %q", string(tc.input), gotHeader, tc.wantHeader)
			}
			if string(gotOut) != string(tc.wantOut) {
				t.Errorf("parseOutputAndBanner(%q) got out %q want out %q", string(tc.input), string(gotOut), string(tc.wantOut))
			}
		})
	}
}

func TestModulesEnv(t *testing.T) {
	// modulesEnv looks at pool.NewGCEConfiguration().BuildEnv().IsProd for
	// special behavior in dev mode. Temporarily override the environment
	// to force testing of the prod configuration.
	old := pool.NewGCEConfiguration().BuildEnv()
	defer pool.NewGCEConfiguration().SetBuildEnv(old)
	pool.NewGCEConfiguration().SetBuildEnv(&buildenv.Environment{
		IsProd: true,
	})

	// In testing we never initialize
	// pool.NewGCEConfiguration().GKENodeHostname(), so we get this odd
	// concatenation back.
	const gkeModuleProxy = "http://:30157"

	testCases := []struct {
		desc string
		st   *buildStatus
		want []string
	}{
		{
			desc: "ec2-builder-repo-non-go",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: "bar"},
				conf: &dashboard.BuildConfig{
					TestHostConf: &dashboard.HostConfig{
						IsReverse: false,
						IsEC2:     true,
					},
				},
			},
			want: []string{"GOPROXY=https://proxy.golang.org"},
		},
		{
			desc: "builder-repo-non-go",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: "bar"},
				conf: &dashboard.BuildConfig{
					TestHostConf: &dashboard.HostConfig{
						IsReverse: false,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=" + gkeModuleProxy},
		},
		{
			desc: "reverse-builder-repo-non-go",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: "bar"},
				conf: &dashboard.BuildConfig{
					TestHostConf: &dashboard.HostConfig{
						IsReverse: true,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=https://proxy.golang.org"},
		},
		{
			desc: "reverse-builder-repo-go",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: ""}, // go
				conf: &dashboard.BuildConfig{
					TestHostConf: &dashboard.HostConfig{
						IsReverse: true,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=off"},
		},
		{
			desc: "builder-repo-go",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: ""}, // go
				conf: &dashboard.BuildConfig{
					TestHostConf: &dashboard.HostConfig{
						IsReverse: false,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=off"},
		},
		{
			desc: "builder-repo-go-outbound-network-allowed",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: ""}, // go
				conf: &dashboard.BuildConfig{
					Name: "test-longtest",
					TestHostConf: &dashboard.HostConfig{
						IsReverse: false,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=" + gkeModuleProxy},
		},
		{
			desc: "reverse-builder-repo-go-outbound-network-allowed",
			st: &buildStatus{
				BuilderRev: buildgo.BuilderRev{SubName: ""}, // go
				conf: &dashboard.BuildConfig{
					Name: "test-longtest",
					TestHostConf: &dashboard.HostConfig{
						IsReverse: true,
						IsEC2:     false,
					},
				},
			},
			want: []string{"GOPROXY=https://proxy.golang.org"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.st.modulesEnv()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("buildStatus.modulesEnv() mismatch (-want, +got)\n%s", diff)
			}
		})
	}
}
