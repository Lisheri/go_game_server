package net

// 处理器
type HandlerFunc func()

// 路由分组
type group struct {
	prefix string;
	handlerMap map[string]HandlerFunc;
}

type router struct {
	group []*group;
}
