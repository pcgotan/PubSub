package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
func TestHelloWorld(t *testing.T) {
	body := gin.H{
		"Hi": "Buddy",
	}
	// Grab our router
	router := Server("0.0.0.0:9000", "/", temp1)
	w := performRequest(router, "POST", "/")

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string

	err := json.Unmarshal([]byte(w.Body.String()), &response)

	value, exists := response["Hi"]

	assert.Nil(t, err)

	assert.True(t, exists)

	assert.Equal(t, body["Hi"], value)
}

func temp1(c *gin.Context) {
	c.JSON(200, gin.H{
		"Hi": "Buddy",
	})
}
