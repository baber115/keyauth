package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "codeup.aliyun.com/baber/go/keyauth/apps/token/api"
	_ "codeup.aliyun.com/baber/go/keyauth/apps/user/api"
)
