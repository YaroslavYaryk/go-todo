package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleRestApi/internal/domain"
	"strconv"
)

type getAllItemsResponse struct {
	Data []domain.TodoItem `json:"data"`
}

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input domain.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.MiddleWare.GetUserById(int(userId))
	fmt.Println(err)
	fmt.Println(user.Id, user.IsPaidMember, "--__---___---_---__---__-----")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if the user is authorized to update this list
	canUserChange, err := h.services.TodoList.IsUserAuthorizedToUpdateList(id, int(user.Id))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't check if user can update this list")
		return
	}

	if !canUserChange {
		newErrorResponse(c, http.StatusForbidden, "User cant update this list, its not his item")
		return
	}

	create, err := h.services.TodoItem.Create(id, input, user)
	if err != nil {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"userId":   userId,
		"todoItem": create,
	})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	items, err := h.services.TodoItem.GetAll(listId, int(userId))

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllItemsResponse{
		Data: items,
	})
}

func (h *Handler) getItemBtId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list_id param")
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item_id param")
		return
	}

	item, err := h.services.TodoItem.GetById(listId, int(userId), itemId)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"Data": item,
	})
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Check if the user is authorized to update this list
	canUserChange, err := h.services.TodoList.IsUserAuthorizedToUpdateList(listId, int(userId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't check if user can update this list")
	}

	if !canUserChange {
		newErrorResponse(c, http.StatusForbidden, "User cant update this list, its not his item")
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item_id param")
		return
	}

	var input domain.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.TodoItem.Update(itemId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "Updated",
	})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	// Check if the user is authorized to update this list
	canUserChange, err := h.services.TodoList.IsUserAuthorizedToUpdateList(listId, int(userId))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Can't check if user can update this list")
	}

	if !canUserChange {
		newErrorResponse(c, http.StatusForbidden, "User cant update this list, its not his item")
		return
	}

	itemId, err := strconv.Atoi(c.Param("item_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item_id param")
		return
	}

	deletedId, err := h.services.TodoItem.Delete(itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "Deleted",
		"id":     deletedId,
	})
}
