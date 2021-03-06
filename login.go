package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUuid returns a randomly generated UUID from Google's UUID library.
func getUuid() string {
	return uuid.New().String()
}

// initCookies use gin-contrib/sessions{/cookie} to initalize a cookie store.
// It generates a random secret for the cookie store -- not ideal for continuity or invalidating previous cookies, but it's secure and it works
func initCookies(r *gin.Engine) {
	r.Use(sessions.Sessions("sarpedon", cookie.NewStore([]byte(getUuid()))))
}

// authRequired provides authentication middleware for ensuring that a user is logged in.
func authRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
	c.Next()
}

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": "Username or password can't be empty."}))
		return
	}

	err := errors.New("Invalid username or password.")

	for _, record := range sarpConfig.Admin {
		if username == record.Username && password == record.Password {
			err = nil
		}
	}

	if err != nil {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": err.Error()}))
		return
	}

	// Save the username in the session
	session.Set("user", username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", pageData(c, "login", gin.H{"error": "Failed to save session."}))
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func getUser(c *gin.Context) string {
	session := sessions.Default(c)
	userName := session.Get("user")
	if userName == nil {
		return ""
	}
	return fmt.Sprintf("%s", userName)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
	session.Delete("user")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}
