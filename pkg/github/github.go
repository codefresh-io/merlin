package github

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/google/go-github/v21/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	DefaultVersion = "latest"
)

type (
	Github interface {
		ReadFile(owner string, repo string, path string, revision string) ([]byte, error)
		GetLatestVersion() string
	}

	git struct {
		client *github.Client
		log    *logrus.Entry
	}
)

func New(token string, log *logrus.Entry) Github {
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
	}
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

func (g *git) GetLatestVersion() string {
	version := DefaultVersion
	releases, _, err := g.client.Repositories.ListReleases(context.Background(), "codefresh-io", "merlin", &github.ListOptions{})
	if err != nil {
		g.log.Errorf("Request to get latest version of venona been rejected , setting version to latest. Original error: %s", err.Error())
		return version
	}
	for _, release := range releases {
		name := strings.Split(*release.Name, "v")
		if len(name) == 2 {
			return name[1]
		}
	}
	return version
}
