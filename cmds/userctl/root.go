package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bufbuild/connect-go"
	"github.com/spf13/cobra"
)

var (
	httpClient connect.HTTPClient = http.DefaultClient
	baseURL    string
	tokenPath  string
)

func getRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use: "userctl [command] [args]",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			content, err := os.ReadFile(tokenPath)
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("failed to read token at %s: %w", tokenPath, err)
				}
			}

			// if we have an access token we wrap the default transport
			// in a custom round-tripper that adds the authentication header
			var rt http.RoundTripper = http.DefaultTransport
			if len(content) > 0 {
				rt = &addHeaderRT{RoundTripper: rt, token: string(content)}
			}

			httpClient = &http.Client{
				Transport: rt,
			}

			return nil
		},
	}

	flags := cmd.PersistentFlags()
	{
		defaultTokenPath := os.Getenv("CIS_IDM_TOKEN_FILE")
		if defaultTokenPath == "" {
			defaultTokenPath = filepath.Join(os.Getenv("HOME"), ".idm-token")
		}

		flags.StringVarP(&baseURL, "url", "U", os.Getenv("CIS_IDM_URL"), "The Base URL for the cis-idm server")
		flags.StringVarP(&tokenPath, "token-file", "t", defaultTokenPath, "The path to the cached access token")
	}

	cmd.AddCommand(
		getLoginCommand(),
		getProfileCommand(),
	)

	return cmd
}

type addHeaderRT struct {
	http.RoundTripper

	token string
}

func (ahr *addHeaderRT) RoundTrip(req *http.Request) (*http.Response, error) {
	if ahr.token != "" {
		req.Header.Add("Authentication", "Bearer "+ahr.token)
	}
	return ahr.RoundTripper.RoundTrip(req)
}
