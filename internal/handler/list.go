package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleRestApi/internal/domain"
	"strconv"
)

type getAllListsResponse struct {
	Data []domain.TodoListExtended `json:"data"`
}

func (h *Handler) createList(c *gin.Context) {

	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input domain.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	create, err := h.services.TodoList.Create(int(userId), input)
	if err != nil {
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"userId":   userId,
		"todoList": create,
	})

}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	lists, err := h.services.TodoList.GetAll(int(userId))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

func (h *Handler) getListBtId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(int(userId), id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Check if the user is authorized to update this list
	canUserChange, err := h.services.TodoList.IsUserAuthorizedToUpdateList(id, int(userId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't check if user can update this list")
	}

	if !canUserChange {
		newErrorResponse(c, http.StatusForbidden, "User cant update this list, its not his item")
		return
	}

	var input domain.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoList.Update(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "Updated",
	})

}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Check if the user is authorized to update this list
	canUserChange, err := h.services.TodoList.IsUserAuthorizedToUpdateList(id, int(userId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't check if user can update this list")
	}

	if !canUserChange {
		newErrorResponse(c, http.StatusForbidden, "User cant update this list, its not his item")
		return
	}

	deletedId, err := h.services.TodoList.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "Deleted",
		"id":     deletedId,
	})

}
