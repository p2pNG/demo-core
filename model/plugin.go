package model

import "github.com/labstack/echo/v4"

type RouterPlugin interface {
	PluginInfo() *PluginInfo
	GetRouter(*echo.Group)
}

type PluginInfo struct {
	Name    string
	Version string
	Prefix  string
}
