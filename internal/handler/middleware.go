package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) userIdentity(c *gin.Context) {

	header := c.Request.Header.Get(authorizationHeader)

	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		newErrorResponse(c, http.StatusUnauthorized, "missing authorization header")
		return
	}

	headerParts := strings.SplitN(header, " ", 2)
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}

	//	parse token
	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set("userId", int64(userId))
}

func getUserId(c *gin.Context) (int64, error) {

	userId, ok := c.Get("userId")

	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "userId not found in context")
		return 0, nil
	}

	idInt, ok2 := userId.(int64)

	if !ok2 {
		newErrorResponse(c, http.StatusInternalServerError, "userId is of invalid type")
	}

	return idInt, nil
}

func isUserAuthorizedToUpdateList(db *sql.DB, userID, listID int) bool {
	var count int
	query := `SELECT COUNT(*) FROM UsersListsTable WHERE user_id = ? AND list_id = ?`
	err := db.QueryRow(query, userID, listID).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
