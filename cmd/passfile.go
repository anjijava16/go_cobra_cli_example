/*
Copyright © 2019 Adron Hall <adron@thrashingcode.com>

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
	"io/ioutil"
	"net/http"
)

// addCmd represents the add command
var passfilecmdapi = &cobra.Command{
	Use:   "passfile",
	Short: "The 'passfile' subcommand will add a passed in key value pair to the application configuration file.",
	Long: `The 'add' subcommand adds a key value pair to the application configuration file. For example:

'<cmd> config add --key theKey --value "the value can be a bunch of things."'.`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, _ := cmd.Flags().GetString("filename")
		fmt.Println(jsonFile)

		// Read the JSON file
		data, err := ioutil.ReadFile(jsonFile)
		if err != nil {
			fmt.Println("Error reading JSON file:", err)
			return
		}

		// Optionally, validate JSON format
		var jsonData map[string]interface{}
		if err := json.Unmarshal(data, &jsonData); err != nil {
			fmt.Println("Invalid JSON format:", err)
			return
		}

		// Send the JSON payload to the API
		apiURL := "https://randomuser.me/api/"
		response, err := sendJSONPayload(apiURL, data)
		if err != nil {
			fmt.Println("Error sending JSON payload:", err)
			return
		}

		fmt.Println("Response from API:", response)
	},
}

// Function to send JSON payload to an API endpoint
func sendJSONPayload(url string, jsonData []byte) (string, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	//req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func init() {
	configCmd.AddCommand(passfilecmdapi)
}
