package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/xamust/couponApp/internal/adapter/api/v1/models"
	"github.com/xamust/couponApp/internal/service/coupon_applier"
	"github.com/xamust/couponApp/internal/usecase"
	customLog "github.com/xamust/couponApp/pkg/logger"
	"log/slog"
	"net/http"
)

type CouponAPIv1 interface {
	CreateCoupon(c echo.Context) error
	FindCouponByID(c echo.Context) error
	CouponList(c echo.Context) error
	DeleteCoupon(c echo.Context) error

	ApplyCoupon(c echo.Context) error
	CouponByUserID(c echo.Context) error
}

type couponAPIv1 struct {
	coupon  usecase.CouponUsecase
	applier coupon_applier.CouponApplier
}

func NewCouponAPIv1(coup usecase.CouponUsecase,
	applier coupon_applier.CouponApplier) CouponAPIv1 {
	return &couponAPIv1{
		coupon:  coup,
		applier: applier,
	}
}

// CreateCoupon
// @Description Метод для создания купона
// @Tags Coupon
// @Accept json
// @Produce json
// @Param requestBody body models.NewAPICoupon true "Request Body, заполнять обязательно"
// @Success 200 {object} models.APICoupon"Success"
// @Router /api/v1/coupon/create [post]
func (co *couponAPIv1) CreateCoupon(c echo.Context) error {
	ctx := c.Request().Context()
	reqModel := &models.NewAPICoupon{}
	if err := c.Bind(&reqModel); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx = customLog.WithCouponName(ctx, reqModel.Name)
	slog.InfoContext(ctx, "CreateCoupon started")
	req, err := reqModel.Map()
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	create, err := co.coupon.Create(ctx, req)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	response := models.MapCouponResp(create)
	slog.InfoContext(ctx, "CreateCoupon done")
	return c.JSON(http.StatusOK, response)
}

// FindCouponByID
// @Description Метод для поиска купона по ID
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} models.APICoupon "Success"
// @Router /api/v1/coupon/:id [get]
func (co *couponAPIv1) FindCouponByID(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	ctx = customLog.WithCouponID(ctx, id)
	slog.InfoContext(ctx, "FindCouponByID started")
	usr, err := co.coupon.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	response := models.MapCouponResp(usr)
	slog.InfoContext(ctx, "FindCouponByID done")
	return c.JSON(http.StatusOK, response)
}

// CouponList
// @Description Метод для вывода списка купонов
// @Tags Coupon
// @Accept json
// @Produce json
// @Param requestBody body models.APICouponList true "Request Body, заполнять обязательно"
// @Success 200 {object} models.APICoupons "Success"
// @Router /api/v1/coupon [post]
func (co *couponAPIv1) CouponList(c echo.Context) error {
	ctx := c.Request().Context()
	reqModel := &models.APICouponList{}
	if err := c.Bind(&reqModel); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	slog.InfoContext(ctx, "CouponList started")
	list, err := co.coupon.List(ctx, reqModel.Limit, reqModel.Offset)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	response := models.MapCouponRespList(list)
	slog.InfoContext(ctx, "CouponList done")
	return c.JSON(http.StatusOK, response)
}

// DeleteCoupon
// @Description Метод для удаления купона по ID
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} models.APICoupon "Success"
// @Router /api/v1/coupon/:id [delete]
func (co *couponAPIv1) DeleteCoupon(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	ctx = customLog.WithCouponID(ctx, id)
	slog.InfoContext(ctx, "DeleteCoupon started")
	usr, err := co.coupon.Delete(ctx, id)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	response := models.MapCouponResp(usr)
	slog.InfoContext(ctx, "DeleteCoupon done")
	return c.JSON(http.StatusOK, response)
}

// ApplyCoupon
// @Description Метод для применения купона к пользователю
// @Tags Coupon
// @Accept json
// @Produce json
// @Param requestBody body models.APICouponApplier true "Request Body, заполнять обязательно"
// @Router /api/v1/coupon/apply [post]
func (co *couponAPIv1) ApplyCoupon(c echo.Context) error {
	ctx := c.Request().Context()
	reqModel := &models.APICouponApplier{}
	if err := c.Bind(&reqModel); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	ctx = customLog.WithCouponID(ctx, reqModel.CouponID)
	ctx = customLog.WithUserID(ctx, reqModel.UserID)
	slog.InfoContext(ctx, "ApplyCoupon started")
	if err := co.applier.Applier(ctx, reqModel.CouponID, reqModel.UserID); err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "ApplyCoupon done")
	return c.NoContent(http.StatusOK)
}

// CouponByUserID
// @Description Метод для поиска купона по UserID
// @Tags Coupon
// @Accept json
// @Produce json
// @Success 200 {object} models.APICoupons "Success"
// @Router /api/v1/coupon/apply/:id [get]
func (co *couponAPIv1) CouponByUserID(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	ctx = customLog.WithUserID(ctx, id)
	slog.InfoContext(ctx, "ApplyCouponByUserID started")
	list, err := co.coupon.ListByUserID(ctx, id, 0, 0)
	if err != nil {
		slog.ErrorContext(customLog.ErrorCtx(ctx, err), "Error: "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	slog.InfoContext(ctx, "ApplyCouponByUserID done")
	return c.JSON(http.StatusOK, models.MapCouponRespList(list))
}
