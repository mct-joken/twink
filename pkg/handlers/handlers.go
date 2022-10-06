package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mct-joken/twink/pkg/docker"
	"io"
	"net/http"
)

var e *echo.Echo
var cli = docker.WorkSpace{}

type containerCreateReq struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	SshPort string `json:"ssh-port"`
}

type containerCreateRes struct {
	ID string `json:"id"`
}

type containerStartRes struct {
	Status string `json:"status"`
}

type containerStopRes struct {
	Status string `json:"status"`
}

type containerDestroyRes struct {
	Status string `json:"status"`
}

func Serve(port string) {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	router()

	docker.NewConnection()
	e.Logger.Fatal(e.Start(":8080"))
}

func router() {
	e.POST("/create", containerCreateHandler)
	e.POST("/container/:id", containerStartHandler)
	e.DELETE("/container/:id", containerDeleteHandler)
	e.DELETE("/container/:id/destroy", containerDestroyHandler)
}

func containerCreateHandler(e echo.Context) error {
	reqBody := containerCreateReq{}
	b, _ := io.ReadAll(e.Request().Body)
	err := json.Unmarshal(b, &reqBody)
	if err != nil {
		fmt.Println(err)
		return e.String(http.StatusBadRequest, "Bad request: BindError")
	}

	create, err := cli.Create(reqBody.Name, reqBody.Image, reqBody.SshPort)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Internal server error")
	}

	marshal, err := json.Marshal(containerCreateRes{ID: create})
	if err != nil {
		return e.String(http.StatusInternalServerError, "Internal server error")
	}

	return e.String(http.StatusCreated, string(marshal))
}

func containerStartHandler(e echo.Context) error {
	id := e.Param("id")
	err := cli.Start(id)
	if err != nil {
		marshal, _ := json.Marshal(containerStartRes{Status: "error"})
		return e.String(http.StatusInternalServerError, string(marshal))
	}

	marshal, err := json.Marshal(containerStartRes{Status: "started"})
	if err != nil {
		return e.String(http.StatusInternalServerError, "Internal server error")
	}

	return e.String(http.StatusOK, string(marshal))
}

func containerDeleteHandler(e echo.Context) error {
	id := e.Param("id")
	err := cli.Stop(id)
	if err != nil {
		marshal, _ := json.Marshal(containerStopRes{Status: "error"})
		return e.String(http.StatusInternalServerError, string(marshal))
	}

	marshal, err := json.Marshal(containerStopRes{Status: "stopped"})
	if err != nil {
		return e.String(http.StatusInternalServerError, "Internal server error")
	}

	return e.String(http.StatusOK, string(marshal))
}

func containerDestroyHandler(e echo.Context) error {
	id := e.Param("id")
	err := cli.Destroy(id)
	if err != nil {
		marshal, _ := json.Marshal(containerDestroyRes{Status: "error"})
		return e.String(http.StatusInternalServerError, string(marshal))
	}

	marshal, err := json.Marshal(containerDestroyRes{Status: "destroyed"})
	if err != nil {
		return e.String(http.StatusInternalServerError, "Internal server error")
	}

	return e.String(http.StatusOK, string(marshal))
}
