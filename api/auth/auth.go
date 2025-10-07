package auth
import (
	"bitresume/api/login"
	"bitresume/config"
	"bitresume/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *gin.Context) {
    url := config.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    c.Redirect(http.StatusTemporaryRedirect, url)
}
// It will be called via goolge auth 
func GoogleCallback(c *gin.Context){
	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token exchange failed"})
		return
	}
	client := config.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed getting user info"})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	json.Unmarshal(body, &userInfo)
	// Step 1: Find user in DB
	user, err := login.GetUserByEmail(userInfo.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not registered"})
		return
	} 
	// Step 2: Create JWT
	jwtToken, err := utils.GenerateJWT(user.Email, user.RollNo, user.Role,user.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	// Step 3: Set cookie
	c.SetCookie("BITRESUME", jwtToken, 3600*6, "/", "localhost", false, true)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head><title>Logged In</title></head>
		<body>
			<script>
				window.opener.postMessage("login-success", "http://localhost:5173");
				window.close();
			</script>
		</body>
		</html>
	`))
}
func Me(c *gin.Context) { // Decode the token and send to frontend
	tokenStr, err := c.Cookie("BITRESUME")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}
	claims, err := utils.ParseJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"email":  claims["email"], 
			"rollNo": claims["rollNo"],
			"role":   claims["role"],  
			"user_name": claims["user_name"],
		},
	})
}
func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie("BITRESUME", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}