/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Response []struct {
	Type      string        `json:"type"`
	Log       []interface{} `json:"log"`
	Msg       []string      `json:"msg"`
	SpamAlias string        `json:"spam_alias"`}

type ErrorResponse struct {
	Type string   `json:"type"`
	Msg  string `json:"msg"`
}

// spamALias represents the alias command
var spamALias = &cobra.Command{
	Use:   "spamAlias",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var username, validity string
		apiKey := viper.GetString("apiKey")
		mxServer := viper.GetString("mxServer")

		addEndpoint := "/add/time_limited_alias"

		if ! viper.GetBool("username") {
			username = viper.GetString("username")
			fmt.Printf("Goto address:\t%s\n", username)
		} else {
			fmt.Println("Please specify a username")
			os.Exit(1)
		}
		validity = viper.GetString("validity")
		fmt.Printf("Valid for:\t%s hour(s)\n", validity)

		requestBody, err := json.Marshal(map[string]string{
			"username": username,
			"validity": validity,
		})

		timeout := time.Duration((10 * time.Second))
		client := http.Client{
			Timeout: timeout,
		}
		request, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/v1%s", mxServer,addEndpoint), bytes.NewBuffer(requestBody))
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("X-API-Key", apiKey)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var checkForError ErrorResponse
		json.Unmarshal(body, &checkForError)

		badResponse := checkForError.Type
		if badResponse != "" {
			fmt.Printf("\n%s: %s\n", checkForError.Type, checkForError.Msg)
		} else {
			var responseObject Response
			json.Unmarshal(body, &responseObject)
			fmt.Printf("\n%s \\O/\t%s\n", responseObject[0].Type, responseObject[0].SpamAlias)
		}
	},
}

func init() {
	createCmd.AddCommand(spamALias)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// spamALias.PersistentFlags().String("foo", "", "A help for foo")
	spamALias.Flags().String("validity", "1", "Validity in hours")
	spamALias.Flags().String("username", "", "Name of a Mailbox to which E-Mails will be forwarded")
	//spamALias.MarkFlagRequired("username")
	viper.BindPFlag("username", spamALias.Flags().Lookup("username"))
	viper.BindPFlag("validity", spamALias.Flags().Lookup("validity"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// spamALias.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
