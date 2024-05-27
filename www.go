package main

import (
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"time"
)

var (
	debug     bool
	tokenFile string
)

type RandomUserResponse struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Email string `json:"email"`
	} `json:"results"`
}

var rootCmd = &cobra.Command{
	Use:   "yourapp",
	Short: "Your application description",
	//Run:   handleToken,
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Handle myapp token",
	Run: func(cmd *cobra.Command, args []string) {
		handleToken(cmd, args)
		fmt.Println("mytest")
		fmt.Println(tokenFile)
	},
}

var myappCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Handle myapp token",
	Run: func(cmd *cobra.Command, args []string) {
		// handleToken(cmd, args)
		fmt.Println("myappCmd")
	},
}

var randomUserNameCMD = &cobra.Command{
	Use:   "randomusername",
	Short: "Handle Random User Name",
	Run: func(cmd *cobra.Command, args []string) {
		handleToken(cmd, args)
		// Make a GET request to the Random User Generator API
		response, err := http.Get("https://randomuser.me/api/")
		if err != nil {
			fmt.Println("Error fetching random user:", err)
			return
		}
		defer response.Body.Close()

		// Decode the JSON response
		var randomUserResponse RandomUserResponse
		if err := json.NewDecoder(response.Body).Decode(&randomUserResponse); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		// Display the random user details
		if len(randomUserResponse.Results) > 0 {
			user := randomUserResponse.Results[0]
			fmt.Println("Random User:")
			fmt.Println("Name:", user.Name.Title, user.Name.First, user.Name.Last)
			fmt.Println("Gender:", user.Gender)
			fmt.Println("Email:", user.Email)
		} else {
			fmt.Println("No results found")
		}
	},
}
var randomUserEmailCmd = &cobra.Command{
	Use:   "randomuseremail",
	Short: "Handle Random User Email",
	Run: func(cmd *cobra.Command, args []string) {
		handleToken(cmd, args)
		fmt.Println("randomuseremail")
	},
}
var randomUserLocationCmd = &cobra.Command{
	Use:   "randomuserlocation",
	Short: "Handle RandomUserLocation",
	Run: func(cmd *cobra.Command, args []string) {
		handleToken(cmd, args)
		fmt.Println("randomuserlocation")
	},
}

var randomuserpictureCmd = &cobra.Command{
	Use:   "randomuserpicture",
	Short: "Handle UserPicture",
	Run: func(cmd *cobra.Command, args []string) {
		handleToken(cmd, args)
		fmt.Println("randomuserpicture")
	},
}

func init() {
	rootCmd.AddCommand(tokenCmd)
	tokenCmd.Flags().StringVar(&tokenFile, "token-file", "", "Path to the token file")
	rootCmd.AddCommand(randomUserNameCMD)
	rootCmd.AddCommand(randomUserEmailCmd)
	rootCmd.AddCommand(randomUserLocationCmd)
	rootCmd.AddCommand(randomuserpictureCmd)

	viper.SetConfigName(".config") // Name of config file (without extension)
	viper.SetConfigType("json")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME")   // Path to look for the config file in
}

func main() {
	fmt.Println(os.Getenv("HOME"))
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleToken(cmd *cobra.Command, args []string) {
	configFilePath := fmt.Sprintf("%s/.config.json", os.Getenv("HOME"))

	// Check if config file exists, if not, create it
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Config file not found; creating a new one.")
		file, err := os.Create(configFilePath)
		if err != nil {
			fmt.Println("Error creating config file:", err)
			return
		}
		file.Close()
	}

	// Attempt to read the config file
	err := viper.ReadInConfig()
	fmt.Println(err)
	//      if err != nil {
	//              fmt.Println("Error reading config file:", err)
	//              return
	//      }

	// Check if token exists
	token := viper.GetString("myapp_token")
	expireStr := viper.GetString("expire_date")
	var expireDate time.Time
	fmt.Println("==========")
	fmt.Println("===========")
	fmt.Println(token)
	fmt.Println("=============")
	fmt.Println("==============")
	if token == "" || isExpired(expireStr) {

		/*
		   if expireStr != "" {
		           expireDate, err = time.Parse(time.RFC3339, expireStr)
		           if err != nil {
		                   fmt.Println("Error parsing expire date:", err)
		                   return
		           }
		   }

		   if token != "" && expireDate.After(time.Now()) {
		           fmt.Println("Token already exists and is valid:", token)
		           return
		   }

		*/
		fmt.Println(token)
		fmt.Println(expireStr)
		// Prompt user for a new token
		prompt := promptui.Prompt{
			Label: "Please enter the new myapp token",
		}

		newToken, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		// Set expiration time (e.g., 1 hour from now)
		expireDate = time.Now().Add(time.Hour)

		// Write the new token and expiration date to the config file
		viper.Set("myapp_token", newToken)
		viper.Set("expire_date", expireDate.Format(time.RFC3339))
		viper.Set("RESOURCE_NAME", "APP-DYNAMIC-myapp-RESOOURE")
		user_name := os.Getenv("USERNAME")
		fmt.Println(user_name)
		username, ok := os.LookupEnv("USERNAME")
		fmt.Println(username)
		if !ok {
			fmt.Println("USERNAME environment variable not found")
			return
		}
		viper.Set("CLIENT_ID", username)
		viper.Set("myapp_URL", "https://myapp-prod.com")

		err = viper.WriteConfig()
		if err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}

		fmt.Println("New token written to config file.")
	} else {
		// Token and expiration date already exist
		fmt.Println("Token and expiration date already exist in config file.")
		fmt.Println("Token:", token)
		fmt.Println("Expiration Date:", expireStr)
	}
}

func isExpired(expireStr string) bool {
	expireTime, err := time.Parse(time.RFC3339, expireStr)
	if err != nil {
		fmt.Println("Error parsing expiration date:", err)
		return true
	}

	// Check if expiration date is within 57 minutes from now
	return time.Until(expireTime) < (57 * time.Minute)
}
