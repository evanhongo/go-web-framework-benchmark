package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	_ "github.com/abemedia/go-don/encoding/text"
	"github.com/astaxie/beego"
	beegoContext "github.com/astaxie/beego/context"
	"github.com/buaazp/fasthttprouter"
	fasthttpSlashRouter "github.com/fasthttp/router"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/julienschmidt/httprouter"
	echo "github.com/labstack/echo/v4"
	echoSlim "github.com/partialize/echo-slim/v4"
	echoSlimMiddleware "github.com/partialize/echo-slim/v4/middleware"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/razonyang/fastrouter"
	"github.com/valyala/fasthttp"
)

var (
	port              = 8080
	sleepTime         = 0
	cpuBound          bool
	target            = 15
	sleepTimeDuration time.Duration
	message           = []byte("hello world")
	messageStr        = "hello world"
	samplingPoint     = 20 // seconds
)

// server [default] [10] [8080]
func main() {
	args := os.Args
	argsLen := len(args)
	webFramework := "default"
	if argsLen > 1 {
		webFramework = args[1]
	}
	if argsLen > 2 {
		sleepTime, _ = strconv.Atoi(args[2])
		if sleepTime == -1 {
			cpuBound = true
			sleepTime = 0
		}
	}
	if argsLen > 3 {
		port, _ = strconv.Atoi(args[3])
	}
	if argsLen > 4 {
		samplingPoint, _ = strconv.Atoi(args[4])
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration := time.Duration(samplingPoint) * time.Second

	go func() {
		time.Sleep(samplingPointDuration)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		var u uint64 = 1024 * 1024
		fmt.Printf("TotalAlloc: %d\n", mem.TotalAlloc/u)
		fmt.Printf("Alloc: %d\n", mem.Alloc/u)
		fmt.Printf("HeapAlloc: %d\n", mem.HeapAlloc/u)
		fmt.Printf("HeapSys: %d\n", mem.HeapSys/u)
	}()

	switch webFramework {
	case "default":
		startDefaultMux()
	case "beego":
		startBeego()
	case "echo":
		startEcho()
	case "echo-slim":
		startEchoSlim()
	case "fasthttp":
		startFasthttp()
	case "fasthttprouter":
		startFastHTTPRouter()
	case "fasthttp/router":
		startFastHTTPSlashRouter()
	case "fasthttp-routing":
		startFastHTTPRouting()
	case "fastrouter":
		startFastRouter()
	case "fiber":
		startFiber()
	case "gin":
		startGin()
	case "httprouter":
		startHTTPRouter()
	default:
		fmt.Println("--------------------------------------------------------------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("--------------------------------------------------------------------")
	}
}

// default mux
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

func startDefaultMux() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

// beego
func beegoHandler(ctx *beegoContext.Context) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	ctx.WriteString(messageStr)
}

func startBeego() {
	beego.BConfig.RunMode = beego.PROD
	beego.BeeLogger.Close()
	mux := beego.NewControllerRegister()
	mux.Get("/hello", beegoHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// echo
func echoHandler(c echo.Context) error {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Response().Write(message)
	return nil
}

func startEcho() {
	e := echo.New()
	e.GET("/hello", echoHandler)

	e.Start(":" + strconv.Itoa(port))
}

// echo-slim
func echoSlimHandler(c echoSlim.Context) error {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Response().Write(message)
	return nil
}

func startEchoSlim() {
	e := echoSlim.New()
	r := echoSlimMiddleware.NewRouter()

	r.GET("/hello", echoSlimHandler)

	e.Use(r.Routes)

	e.Start(":" + strconv.Itoa(port))
}

func startFasthttp() {
	s := &fasthttp.Server{
		Handler:                       fastHTTPHandler,
		GetOnly:                       true,
		NoDefaultDate:                 true,
		NoDefaultContentType:          true,
		DisableHeaderNamesNormalizing: true,
	}

	log.Fatal(s.ListenAndServe(":" + strconv.Itoa(port)))
}

// fasthttprouter
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}

	ctx.Write(message)
}

func startFastHTTPRouter() {
	mux := fasthttprouter.New()
	mux.GET("/hello", fastHTTPHandler)
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.Handler)
}

// fasthttp Router
func startFastHTTPSlashRouter() {
	mux := fasthttpSlashRouter.New()
	mux.GET("/hello", fastHTTPHandler)
	log.Fatal(fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.Handler))
}

// fasthttprouting
func fastHTTPRoutingHandler(c *routing.Context) error {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Write(message)
	return nil
}

func startFastHTTPRouting() {
	mux := routing.New()
	mux.Get("/hello", fastHTTPRoutingHandler)
	fasthttp.ListenAndServe(":"+strconv.Itoa(port), mux.HandleRequest)
}

// fastrouter
func fastRouterHandler(w http.ResponseWriter, r *http.Request) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

func startFastRouter() {
	mux := fastrouter.New()
	mux.Get("/hello", fastRouterHandler)
	mux.Prepare()
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// fiber
func fiberHandler(c *fiber.Ctx) error {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	return c.SendString(messageStr)
}

func startFiber() {
	app := fiber.New(fiber.Config{
		Prefork:                   true,
		CaseSensitive:             true,
		StrictRouting:             true,
		DisableDefaultDate:        true,
		DisableHeaderNormalizing:  true,
		DisableDefaultContentType: true,
	})
	app.Get("/hello", fiberHandler)
	log.Fatal(app.Listen(":" + strconv.Itoa(port)))
}

// gin
func ginHandler(c *gin.Context) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Writer.Write(message)
}

func startGin() {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.GET("/hello", ginHandler)
	mux.Run(":" + strconv.Itoa(port))
}

// httprouter
func httpRouterHandler(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

func startHTTPRouter() {
	mux := httprouter.New()
	mux.GET("/hello", httpRouterHandler)
	http.ListenAndServe(":"+strconv.Itoa(port), mux)
}

// mock
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}
