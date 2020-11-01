package core

import (
	"git.ixarea.com/p2pNG/p2pNG-core/model"
)

var plugins map[string]model.RouterPlugin

func init() {
	plugins = make(map[string]model.RouterPlugin)
}

func RegisterRouterPlugin(plugin model.RouterPlugin) {
	x := plugin.PluginInfo()
	plugins[x.Name] = plugin
}
func GetRouterPluginRegistry() []model.RouterPlugin {
	x := make([]model.RouterPlugin, len(plugins))
	i := 0
	for name := range plugins {
		x[i] = plugins[name]
		i++
	}
	return x
}

func GetRouterPlugin(name string) (model.RouterPlugin, bool) {
	x, ok := plugins[name]
	return x, ok
}
