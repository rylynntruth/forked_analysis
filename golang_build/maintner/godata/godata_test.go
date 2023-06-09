// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package godata

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/build/gerrit"
	"golang.org/x/build/internal/secret"
	"golang.org/x/build/maintner"
)

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := Get(context.Background())
		if err != nil {
			b.Fatal(err)
		}
	}
}

var (
	corpusMu    sync.Mutex
	corpusCache *maintner.Corpus
)

func getGoData(tb testing.TB) *maintner.Corpus {
	if testing.Short() {
		tb.Skip("skipping test requiring large download in short mode")
	}
	corpusMu.Lock()
	defer corpusMu.Unlock()
	if corpusCache != nil {
		return corpusCache
	}
	var err error
	corpusCache, err = Get(context.Background())
	if err != nil {
		tb.Fatalf("getting corpus: %v", err)
	}
	return corpusCache
}

func TestCorpusCheck(t *testing.T) {
	c := getGoData(t)
	if err := c.Check(); err != nil {
		t.Fatal(err)
	}
}

func TestGerritForeachNonChangeRef(t *testing.T) {
	c := getGoData(t)
	c.Gerrit().ForeachProjectUnsorted(func(gp *maintner.GerritProject) error {
		t.Logf("%s:", gp.ServerSlashProject())
		gp.ForeachNonChangeRef(func(ref string, hash maintner.GitHash) error {
			t.Logf("  %s %s", hash, ref)
			return nil
		})
		return nil
	})
}

// In the past, some Gerrit ref changes came before the git in the log.
// This tests that we handle Gerrit meta changes that happen before
// the referenced git commit is known.
func TestGerritOutOfOrderMetaChanges(t *testing.T) {
	c := getGoData(t)

	// Merged:
	goProj := c.Gerrit().Project("go.googlesource.com", "go")
	cl := goProj.CL(38634)
	if cl == nil {
		t.Fatal("CL 38634 not found")
	}
	if g, w := cl.Status, "merged"; g != w {
		t.Errorf("CL status = %q; want %q", g, w)
	}

	// Deleted:
	gddo := c.Gerrit().Project("go.googlesource.com", "gddo")
	cl = gddo.CL(37452)
	if cl == nil {
		t.Fatal("CL 37452 not found")
	}
	t.Logf("Got: %+v", *cl)
}

func TestGerritSkipPrivateCLs(t *testing.T) {
	c := getGoData(t)
	proj := c.Gerrit().Project("go.googlesource.com", "gddo")
	proj.ForeachOpenCL(func(cl *maintner.GerritCL) error {
		if cl.Number == 37452 {
			t.Error("unexpected private CL 37452")
		}
		return nil
	})
}

func TestGerritMetaNonNil(t *testing.T) {
	c := getGoData(t)
	c.Gerrit().ForeachProjectUnsorted(func(gp *maintner.GerritProject) error {
		var maxCL int32
		gp.ForeachCLUnsorted(func(cl *maintner.GerritCL) error {
			if cl.Meta == nil {
				t.Errorf("%s: ForeachCLUnsorted-enumerated CL %d has nil Meta", gp.ServerSlashProject(), cl.Number)
			}
			if len(cl.Metas) == 0 {
				t.Errorf("%s: ForeachCLUnsorted-enumerated CL %d has empty Metas", gp.ServerSlashProject(), cl.Number)
			}
			if cl.Commit == nil {
				t.Errorf("%s: ForeachCLUnsorted-enumerated CL %d has nil Commit", gp.ServerSlashProject(), cl.Number)
			}
			if cl.Number > maxCL {
				maxCL = cl.Number
			}
			return nil
		})
		gp.ForeachOpenCL(func(cl *maintner.GerritCL) error {
			if cl.Meta == nil {
				t.Errorf("%s: ForeachOpenCL-enumerated CL %d has nil Meta", gp.ServerSlashProject(), cl.Number)
			}
			if len(cl.Metas) == 0 {
				t.Errorf("%s: ForeachOpenCL-enumerated CL %d has empty Metas", gp.ServerSlashProject(), cl.Number)
			}
			if cl.Commit == nil {
				t.Errorf("%s: ForeachOpenCL-enumerated CL %d has nil Commit", gp.ServerSlashProject(), cl.Number)
			}
			if cl.Number > maxCL {
				t.Fatalf("%s: ForeachOpenCL-enumerated CL %d higher than max CL %d from ForeachCLUnsorted", gp.ServerSlashProject(), cl.Number, maxCL)
			}
			return nil
		})

		// And test that CL won't yield an incomplete one either:
		for n := int32(0); n <= maxCL; n++ {
			cl := gp.CL(n)
			if cl == nil {
				continue
			}
			if cl.Meta == nil {
				t.Errorf("%s: CL(%d) has nil Meta", gp.ServerSlashProject(), cl.Number)
			}
			if len(cl.Metas) == 0 {
				t.Errorf("%s: CL(%d) has empty Metas", gp.ServerSlashProject(), cl.Number)
			}
			if cl.Commit == nil {
				t.Errorf("%s: CL(%d) has nil Commit", gp.ServerSlashProject(), cl.Number)
			}
		}
		return nil
	})
}

func TestGitAncestor(t *testing.T) {
	c := getGoData(t)
	tests := []struct {
		subject, ancestor string
		want              bool
	}{
		{"3b5637ff2bd5c03479780995e7a35c48222157c1", "0bb0b61d6a85b2a1a33dcbc418089656f2754d32", true},
		{"0bb0b61d6a85b2a1a33dcbc418089656f2754d32", "3b5637ff2bd5c03479780995e7a35c48222157c1", false},

		{"8f06e217eac10bae4993ca371ade35fecd26270e", "22f1b56dab29d397d2bdbdd603d85e60fb678089", true},
		{"22f1b56dab29d397d2bdbdd603d85e60fb678089", "8f06e217eac10bae4993ca371ade35fecd26270e", false},

		// Was crashing. Issue 22753.
		{"3a181dc7bc8fd0c61d6090a85f87c934f1874802", "f65abf6ddc8d1f3d403a9195fd74eaffa022b07f", true},
		// The reverse of the above, to try to reproduce the
		// panic if I got the order backwards:
		{"f65abf6ddc8d1f3d403a9195fd74eaffa022b07f", "3a181dc7bc8fd0c61d6090a85f87c934f1874802", false},

		// Same on both sides:
		{"0bb0b61d6a85b2a1a33dcbc418089656f2754d32", "0bb0b61d6a85b2a1a33dcbc418089656f2754d32", false},
		{"3b5637ff2bd5c03479780995e7a35c48222157c1", "3b5637ff2bd5c03479780995e7a35c48222157c1", false},
	}
	for i, tt := range tests {
		subject := c.GitCommit(tt.subject)
		if subject == nil {
			t.Errorf("%d. missing subject commit %q", i, tt.subject)
			continue
		}
		anc := c.GitCommit(tt.ancestor)
		if anc == nil {
			t.Errorf("%d. missing ancestor commit %q", i, tt.ancestor)
			continue
		}
		got := subject.HasAncestor(anc)
		if got != tt.want {
			t.Errorf("HasAncestor(%q, %q) = %v; want %v", tt.subject, tt.ancestor, got, tt.want)
		}
	}
}

func BenchmarkGitAncestor(b *testing.B) {
	c := getGoData(b)
	subject := c.GitCommit("3b5637ff2bd5c03479780995e7a35c48222157c1")
	anc := c.GitCommit("0bb0b61d6a85b2a1a33dcbc418089656f2754d32")
	if subject == nil || anc == nil {
		b.Fatal("missing commit(s)")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !subject.HasAncestor(anc) {
			b.Fatal("wrong answer")
		}
	}
}

// Issue 23007: a Gerrit CL can switch branches. Make sure we handle that.
func TestGerritCLChangingBranches(t *testing.T) {
	c := getGoData(t)

	tests := []struct {
		server, project string
		cl              int32
		want            string
	}{
		// Changed branch in the middle:
		// (Unsubmitted at the time of this test, so if it changes back, this test
		// may break.)
		{"go.googlesource.com", "go", 33776, "master"},

		// Submitted to boringcrypto:
		{"go.googlesource.com", "go", 82138, "dev.boringcrypto"},
		// Submitted to master:
		{"go.googlesource.com", "go", 83578, "master"},
	}

	for _, tt := range tests {
		cl := c.Gerrit().Project(tt.server, tt.project).CL(tt.cl)
		if got := cl.Branch(); got != tt.want {
			t.Errorf("%q, %q, CL %d = branch %q; want %q", tt.server, tt.project, tt.cl, got, tt.want)
		}
	}
}

func TestGerritHashTags(t *testing.T) {
	c := getGoData(t)
	cl := c.Gerrit().Project("go.googlesource.com", "go").CL(81778)
	want := `added "bar, foo" = "bar,foo"
removed "bar" = "foo"
removed "foo" = ""
added "bar, foo" = "bar,foo"
removed "bar" = "foo"
added "bar" = "bar,foo"
added "blarf, quux" removed "foo" = "bar,quux,blarf"
removed "bar" = "quux,blarf"
`

	var log bytes.Buffer
	for _, meta := range cl.Metas {
		added, removed, ok := meta.HashtagEdits()
		if ok {
			if added != "" {
				fmt.Fprintf(&log, "added %q ", added)
			}
			if removed != "" {
				fmt.Fprintf(&log, "removed %q ", removed)
			}
			fmt.Fprintf(&log, "= %q\n", meta.Hashtags())
		}
	}
	got := log.String()
	if !strings.HasPrefix(got, want) {
		t.Errorf("got:\n%s\n\nwant prefix:\n%s", got, want)
	}
}

// getSecret retrieves a secret by name from the secret manager service.
func getSecret(name string) (string, error) {
	sc, err := secret.NewClient()
	if err != nil {
		return "", err
	}
	defer sc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return sc.Retrieve(ctx, name)
}

func getGerritAuth() (username string, password string, err error) {
	var slurp string
	if metadata.OnGCE() {
		slurp, _ = getSecret(secret.NameGobotPassword)
	}
	if slurp == "" {
		var ok bool
		slurp, ok = os.LookupEnv("TEST_GERRIT_AUTH")
		if !ok {
			return "", "", errors.New("environment variable TEST_GERRIT_AUTH is not set")
		}
	}
	f := strings.SplitN(strings.TrimSpace(slurp), ":", 2)
	if len(f) == 1 {
		// assume the whole thing is the token
		return "git-gobot.golang.org", f[0], nil
	}
	if len(f) != 2 || f[0] == "" || f[1] == "" {
		return "", "", fmt.Errorf("Expected Gerrit token %q to be of form <git-email>:<token>", slurp)
	}
	return f[0], f[1], nil
}

// Hit the Gerrit API and compare its computation of CLs' hashtags against what maintner thinks.
// Off by default unless $TEST_GERRIT_AUTH is defined with "user:token", or we're running in the
// prod project.
func TestGerritHashtags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	c := getGoData(t)
	user, pass, err := getGerritAuth()
	if err != nil {
		t.Skipf("no Gerrit auth defined, skipping: %v", err)
	}
	gc := gerrit.NewClient("https://go-review.googlesource.com", gerrit.BasicAuth(user, pass))
	ctx := context.Background()
	more := true
	n := 0
	for more {
		// We search Gerrit for "hashtag", which seems to also
		// search auto-generated gerrit meta (notedb) texts,
		// so this has the effect of searching for all Gerrit
		// changes that have ever had hashtags added or
		// removed:
		cis, err := gc.QueryChanges(ctx, "hashtag", gerrit.QueryChangesOpt{
			Start: n,
		})
		if err != nil {
			t.Fatal(err)
		}
		for _, ci := range cis {
			n++
			cl := c.Gerrit().Project("go.googlesource.com", ci.Project).CL(int32(ci.ChangeNumber))
			if cl == nil {
				t.Logf("Ignoring not-in-maintner %s/%v", ci.Project, ci.ChangeNumber)
				continue
			}
			sort.Strings(ci.Hashtags)
			want := strings.Join(ci.Hashtags, ", ")
			got := canonicalTagList(string(cl.Meta.Hashtags()))
			if got != want {
				t.Errorf("ci: https://golang.org/cl/%d (%s) -- maintner = %q; want gerrit value %q", ci.ChangeNumber, ci.Project, got, want)
			}
			more = ci.MoreChanges
		}
	}
	t.Logf("N = %v", n)
}

func canonicalTagList(s string) string {
	var sl []string
	for _, v := range strings.Split(s, ",") {
		sl = append(sl, strings.TrimSpace(v))
	}
	sort.Strings(sl)
	return strings.Join(sl, ", ")
}
