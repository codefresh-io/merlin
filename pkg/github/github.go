package github

import (
	"bytes"
	"context"
	"net/http"

	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type (
	Github interface {
		ReadFile(owner string, repo string, path string, revision string) ([]byte, error)
	}

	git struct {
		client *github.Client
		log    *logrus.Entry
	}
)

func New(token string, log *logrus.Entry) (Github, error) {
	var tc *http.Client

	if token != "" {
		log.Debug("Using token")
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)
	return &git{
		client: client,
		log:    log,
	}, nil
}

func (g *git) ReadFile(owner string, repo string, path string, revision string) ([]byte, error) {
	ctx := context.Background()
	reader, err := g.client.Repositories.DownloadContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
		Ref: revision,
	})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.Bytes(), nil
}
