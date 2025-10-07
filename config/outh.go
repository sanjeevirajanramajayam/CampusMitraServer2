package config
import (
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig *oauth2.Config
func InitOAuth() {
	requiredVars := []string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET", "GOOGLE_REDIRECT_URL"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			log.Fatalf("Missing required environment variable: %s", v)
		}
	}
	GoogleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	fmt.Printf("OAuth Config Loaded:\nClientID: %s\nRedirectURL: %s\n", 
		GoogleOauthConfig.ClientID, 
		GoogleOauthConfig.RedirectURL)
}
