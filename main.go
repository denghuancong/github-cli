package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/cobra"
)

var (
	token        string
	slackUsed    bool
	slackToken   string
	slackChannel string
	configFile   string
	globalCtx    context.Context
	globalCfg    *Config
	globalClient *Client
	output       bytes.Buffer
)

func perror(err error) {
	if err == nil {
		return
	}

	println(err.Error())
	os.Exit(1)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "github-cli",
		Short: "Github CLI",
	}

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "C", "", "Config File, default ~/.github-cli/config.toml")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "T", "", "Github Token")
	rootCmd.PersistentFlags().BoolVar(&slackUsed, "slack", false, "Send result to slack")
	rootCmd.PersistentFlags().StringVar(&slackToken, "slack_token", "", "Slack Token")
	rootCmd.PersistentFlags().StringVar(&slackChannel, "slack_channel", "", "Slack channel")

	rootCmd.AddCommand(
		newPullsCommand(),
		newPullCommand(),
		newIssuesCommand(),
		newIssueCommand(),
		newTrendingCommand(),
		newEventsCommand(),
	)

	cobra.OnInitialize(initGlobal)
	cobra.EnablePrefixMatching = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(rootCmd.UsageString())
	}

	fmt.Println(output.String())

	if slackUsed {
		SendToSlack(globalCfg.Slack, output.String())
	}
}

func initGlobal() {
	usr, err := user.Current()
	perror(err)

	if len(configFile) == 0 {
		configFile = path.Join(usr.HomeDir, ".github-cli/config.toml")
	}
	cfg, err := NewConfigFromFile(configFile)
	perror(err)

	if len(token) > 0 {
		cfg.Token = token
	}

	if len(slackToken) > 0 {
		cfg.Slack.Token = slackToken
	}

	if len(slackChannel) > 0 {
		cfg.Slack.Channel = slackChannel
	}

	globalCtx = context.Background()
	globalCfg = cfg
	globalClient = NewClient(globalCtx, cfg)
}
