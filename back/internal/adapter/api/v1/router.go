package v1

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xamust/couponApp/internal/service/coupon_applier"
	"github.com/xamust/couponApp/internal/usecase"
	"net/http"
)

type Router struct {
	user    *usecase.UserUsecase
	coup    *usecase.CouponUsecase
	applier *coupon_applier.CouponApplier
}

func NewRouter(
	user *usecase.UserUsecase,
	coup *usecase.CouponUsecase,
	applier *coupon_applier.CouponApplier,
) *Router {
	return &Router{
		user:    user,
		coup:    coup,
		applier: applier,
	}
}

func (r *Router) Build(ctx context.Context) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true

	// v1
	v1 := e.Group("/api/v1",
		middleware.Logger(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://front:3000"},
			AllowMethods: []string{"GET", "POST", "DELETE"},
			AllowHeaders: []string{"Content-Type"},
		}))

	// user
	u := v1.Group("/user")
	usrAPIv1 := NewUserAPIv1(*r.user)
	u.POST("/create", usrAPIv1.CreateUser)
	u.GET("/:id", usrAPIv1.FindUserByID)
	u.POST("", usrAPIv1.UsersList)
	u.DELETE("/:id", usrAPIv1.DeleteUser)

	// coupon
	c := v1.Group("/coupon")
	coupAPIv1 := NewCouponAPIv1(*r.coup, *r.applier)
	c.POST("/create", coupAPIv1.CreateCoupon)
	c.GET("/:id", coupAPIv1.FindCouponByID)
	c.POST("", coupAPIv1.CouponList)
	c.DELETE("/:id", coupAPIv1.DeleteCoupon)

	// apply
	c.POST("/apply", coupAPIv1.ApplyCoupon)
	c.GET("/apply/:id", coupAPIv1.CouponByUserID)

	//default routes
	e.Group("/health").GET("*", healthHandler)
	e.Group("/admin").Any("*", echoSwagger.WrapHandler)

	return e, nil
}

func healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
