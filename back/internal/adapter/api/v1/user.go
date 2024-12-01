package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/xamust/couponApp/internal/adapter/api/v1/models"
	"github.com/xamust/couponApp/internal/usecase"
	customLog "github.com/xamust/couponApp/pkg/logger"
	"golang.org/x/exp/slog"
	"net/http"
)

type UserAPIv1 interface {
	CreateUser(c echo.Context) error
	FindUserByID(c echo.Context) error
	UsersList(c echo.Context) error
	DeleteUser(c echo.Context) error
}

type userAPIv1 struct {
	user usecase.UserUsecase
}

func NewUserAPIv1(user usecase.UserUsecase) UserAPIv1 {
	return &userAPIv1{
		user: user,
	}
}

// CreateUser
// @Description Метод для создания пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param requestBody body models.NewAPIUser true "Request Body, заполнять обязательно"
// @Success 200 {object} models.APIUser "Success"
// @Router /api/v1/user/create [post]
func (u *userAPIv1) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()
	reqModel := &models.NewAPIUser{}
	if err := c.Bind(&reqModel); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx = customLog.WithUserName(ctx, reqModel.Name)
	slog.InfoContext(ctx, "CreateUser started")
	create, err := u.user.Create(ctx, reqModel.Map())
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "CreateUser done")
	return c.JSON(http.StatusOK, models.MapUserResp(create))
}

// FindUserByID
// @Description Метод для поиска пользователя по ID
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.APIUser "Success"
// @Router /api/v1/user/:id [get]
func (u *userAPIv1) FindUserByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	ctx = customLog.WithUserID(ctx, id)
	slog.InfoContext(ctx, "FindUserByID started")
	usr, err := u.user.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "FindUserByID done")
	return c.JSON(http.StatusOK, models.MapUserResp(usr))
}

// UsersList
// @Description Метод для вывода списка пользователей
// @Tags User
// @Accept json
// @Produce json
// @Param requestBody body models.APIUserList true "Request Body, заполнять обязательно"
// @Success 200 {object} models.APIUsers "Success"
// @Router /api/v1/user [post]
func (u *userAPIv1) UsersList(c echo.Context) error {
	ctx := c.Request().Context()
	reqModel := &models.APIUserList{}
	if err := c.Bind(&reqModel); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	slog.InfoContext(ctx, "UsersList started")
	list, err := u.user.List(ctx, reqModel.Limit, reqModel.Offset)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "UsersList done")
	return c.JSON(http.StatusOK, models.MapUserRespList(list))
}

// DeleteUser
// @Description Метод для удаления пользователя по ID
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.APIUser "Success"
// @Router /api/v1/user/:id [delete]
func (u *userAPIv1) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	ctx = customLog.WithUserID(ctx, id)
	slog.InfoContext(ctx, "DeleteUser started")
	usr, err := u.user.Delete(ctx, id)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "DeleteUser done")
	return c.JSON(http.StatusOK, models.MapUserResp(usr))
}
