// Author: Steve Zhang
// Date: 2020/9/23 3:11 下午

package main

import (
	"fmt"

	"go-server/component"
)

// setupComponent 配置组件
func setupComponent(conf string, port int) (err error) {

	// 配置配置组件
	if err = component.SetupConf(conf); err != nil {
		err = fmt.Errorf("component.SetupConf(%s): %w", conf, err)
		return
	}

	// 配置消息日志组件
	if err = component.SetupInfLogger(); err != nil {
		err = fmt.Errorf("component.SetupInfLogger: %w", err)
		return
	}

	// 配置错误日志组件
	if err = component.SetupErrLogger(); err != nil {
		err = fmt.Errorf("component.SetupErrLogger: %w", err)
		return
	}

	// 配置缓存
	//if err = component.SetupCache(); err != nil {
	//	err = fmt.Errorf("component.SetupCache: %w", err)
	//	return
	//}

	// 配置DB
	//if err = component.SetupDB(); err != nil {
	//	err = fmt.Errorf("component.SetupDB: %w", err)
	//	return
	//}

	// 配置HTTP服务
	if err = component.SetupHttpServer(port); err != nil {
		err = fmt.Errorf("component.SetupHttpServer(%d): %v", port, err)
		return
	}

	return
}
