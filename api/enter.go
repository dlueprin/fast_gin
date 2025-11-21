package api

import "fast_gin/api/user_api"

// api总入口，通过点api再点其他的来实现统一管理
type Api struct {
	UserApi user_api.UserApi
}

// 全局变量，全局单例模式，延迟初始化，放在global会导致增加不必要的耦合，而且这个上层变量可以依赖下层基础变量，而不能反过来依赖
var App = new(Api)
