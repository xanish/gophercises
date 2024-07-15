package sitemap_builder

import (
	"bytes"
	approvals "github.com/approvals/go-approval-tests"
	"testing"
)

func TestBuildSitemap(t *testing.T) {
	got := bytes.Buffer{}
	urls := []string{
		"https://gophercises.com/demos/cyoa",
		"https://gophercises.com",
	}

	BuildSitemap(urls, &got)

	approvals.VerifyString(t, got.String())
}
