package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/joho/godotenv"
	"github.com/kemboy-svg/investment/routes"
	"github.com/kemboy-svg/investment/store"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or failed to load")
	}
	if err := store.DbInit(); err != nil {
		fmt.Println(err, "=DB init errors=")
	}
	st := store.Store{}
	st.MigrateAllModels()
	// controllers.DeploymentFunctions(
}

func main() {
	e := echo.New()
	routes.Routes(e)

	// godotenv.Load()
	var port string
	if os.Getenv("IS_PRODUCTION") == "TRUE" {
		port = os.Getenv("PORT")
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			fmt.Println("Req: ", string(reqBody))
			fmt.Println("Resp: ", string(resBody))
		}))
	} else {
		port = os.Getenv("PORT_DEV")
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			fmt.Println("Req: ", string(reqBody))
			fmt.Println("Resp: ", string(resBody))
		}))
	}
	if port == "" {
		port = "8450"
	}
	// port = "8500"
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("20M"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, "Welcome to my APIS portal")
	})
	
	config := middleware.RateLimiterConfig{

		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 5, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

	// =====> PROMETHEUS
	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		// labels of default metrics can be modified or added with `LabelFuncs` function
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"scheme": func(c echo.Context, err error) string { // additional custom label
				return c.Scheme()
			},
			"host": func(c echo.Context, err error) string { // overrides default 'host' label value
				return "y_" + c.Request().Host
			},
		},
	
		HistogramOptsFunc: func(opts prometheus.HistogramOpts) prometheus.HistogramOpts {
			if opts.Name == "request_duration_seconds" {
				opts.Buckets = []float64{1000.0, 10_000.0, 100_000.0, 1_000_000.0} // 1KB ,10KB, 100KB, 1MB
			}
			return opts
		},
		CounterOptsFunc: func(opts prometheus.CounterOpts) prometheus.CounterOpts {
			if opts.Name == "requests_total" {
				opts.ConstLabels = prometheus.Labels{"my_const": "123"}
			}
			return opts
		},
	})) // adds middleware to gather metrics
	e.GET("/metrics_prometheus", echoprometheus.NewHandler())

	e.Logger.Fatal(e.Start(":" + port))

}

// momato@10
// 20222023$
