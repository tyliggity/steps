package base

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andygrunwald/go-jira"
	"github.com/stackpulse/steps-sdk-go/log"
	"golang.org/x/oauth2"
)

type AtlasssianResource struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes,omitempty"`
	AvatarURL string   `json:"avatarUrl,omitempty"`
}

const (
	atlassianListResourcesURL = "https://api.atlassian.com/oauth/token/accessible-resources"
	jiraOauthAPIURLTemplate   = "https://api.atlassian.com/ex/jira/%s"
)

type Args struct {
	SiteName  string `env:"JIRA_SITE_NAME,required"`
	JiraToken string `env:"JIRA_TOKEN,required"`
}

func GetOauthConnectionDetails(token *oauth2.Token, site string) (client *http.Client, url string, err error) {
	client = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))

	request, err := http.NewRequest(http.MethodGet, atlassianListResourcesURL, nil)
	if err != nil {
		return client, url, err
	}

	request.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return client, url, err
	}
	log.Debug("get response code [%s]", resp.Status)

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return client, url, fmt.Errorf("a non-success response code was received [%s], access token might have expired or Jira API may be down", resp.Status)
	}

	defer resp.Body.Close()
	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return client, url, err
	}
	log.Debug("get response output [%s]", string(output))

	var resources []AtlasssianResource
	err = json.Unmarshal(output, &resources)
	if err != nil {
		return client, "", err
	}

	for _, resource := range resources {
		if resource.Name == site {
			url = fmt.Sprintf(jiraOauthAPIURLTemplate, resource.ID)
			return client, url, nil
		}
	}

	return client, url, fmt.Errorf("failed to find site specified in JIRA_SITE_NAME in the sites accessible by the token")
}

func NewClient(token *oauth2.Token, site string) (*jira.Client, error) {
	oauthClient, url, err := GetOauthConnectionDetails(token, site)
	if err != nil {
		return nil, fmt.Errorf("get oauth details: %w", err)
	}

	jiraClient, err := jira.NewClient(oauthClient, url)
	if err != nil {
		return nil, fmt.Errorf("initialize client: %w", err)
	}

	return jiraClient, nil
}

func ExtractResponse(resp *jira.Response) []byte {
	if resp == nil {
		return nil
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return content
}
