package category_http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	category_application "github.com/dhanielsales/golang-scaffold/modules/category/application"
)

type CategoryController struct {
	service *category_application.CategoryService
}

func NewCategoryController(service *category_application.CategoryService) *CategoryController {
	return &CategoryController{
		service: service,
	}
}

type createCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=50"`
	Description string `json:"description" validate:"max=300"`
}

// @Summary Create one category.
// @Description creates one category.
// @Tags Category
// @Accept */*
// @Produce json
// @Param Category body createCategoryRequest true "Category to create"
// @Success 201
// @Router /category [post]
func (t *CategoryController) create(c echo.Context) error {
	var req createCategoryRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request body",
		})
	}

	err := t.service.Create(c.Request().Context(), category_application.CreateCategoryPayload{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to create category",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary Get all categories.
// @Description fetch every category available.
// @Tags Category
// @Accept */*
// @Produce json
// @Success 200 {object} []category_entity.Category
// @Router /category [get]
func (t *CategoryController) getAll(c echo.Context) error {
	categories, err := t.service.GetAll(c.Request().Context())
	if err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to get category",
		})
	}

	return c.JSON(http.StatusOK, categories)
}
