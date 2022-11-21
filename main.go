package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
)

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

func main() {
	key := os.Getenv("token")
	httpClient := http.Client{
		Transport: &authedTransport{
			key:     key,
			wrapped: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)
	var viewerResp *getViewerResponse
	ctx := context.Background()
	viewerResp, err := getViewer(ctx, graphqlClient)
	if err != nil {
		return
	}
	fmt.Println("you are", viewerResp.Viewer.Login, "created on", viewerResp.Viewer.CreatedAt.Format("2006-01-02"))

	repoResp, err := countRepoIssuesInitial(context.Background(), graphqlClient)
	if err != nil {
		fmt.Println("", err)
	}

	var issuesurls []string

	for _, v := range repoResp.Repository.Issues.Edges {
		issuesurls = append(issuesurls, v.Node.Url)
	}
