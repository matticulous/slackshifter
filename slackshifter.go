package main

import (
	"fmt"
	"os"
	"strings"

	slack "github.com/huguesalary/slack-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Package-level vars for flag values
var (
	url      string
	channel  string
	username string
	alias    string
	iconURL  string
)

var rootCmd = &cobra.Command{
	Use:  "slackshifter <text to send>",
	Args: cobra.MinimumNArgs(1),
	RunE: shift,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&url, "url", "", "Slack webhook URL") //Required
	rootCmd.Flags().StringVarP(&channel, "channel", "c", "", "Slack channel to post message")
	rootCmd.Flags().StringVar(&username, "username", "", "Slack username to post as")
	rootCmd.Flags().StringVarP(&iconURL, "icon-url", "i", "", "URL for icon")
	rootCmd.Flags().StringVarP(&alias, "alias", "a", "", "Alias config to use")
	viper.BindPFlags(rootCmd.Flags())
	viper.BindPFlag("webhook_url", rootCmd.Flags().Lookup("url"))
	viper.BindPFlag("icon_url", rootCmd.Flags().Lookup("icon-url"))
}

func shift(cmd *cobra.Command, args []string) error {
	if aliasName := viper.GetString("alias"); aliasName != "" {
		alias := viper.GetStringMap(aliasName)

		username = alias["username"].(string)
		iconURL = alias["icon_url"].(string)
	}
	text := strings.Join(args[0:], " ")

	s := slack.NewClient(viper.GetString("webhook_url"))
	m := &slack.Message{
		Text: text,
	}
	if username != "" {
		m.Username = username
	}
	if channel != "" {
		m.Channel = channel
	}
	if iconURL != "" {
		m.IconUrl = iconURL
	}

	err := s.SendMessage(m)
	if err != nil {
		return err
	}
	return nil
}

func initConfig() {
	viper.SetEnvPrefix("slackshifter")
	viper.AutomaticEnv()                // read in environment variables that match
	viper.SetConfigName("slackshifter") // name of config file (without extension)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
