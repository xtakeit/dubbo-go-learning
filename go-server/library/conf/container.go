package conf

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

// CompareObjConfRst 定义对比对象配置结果类型
type CompareObjConfRst uint32

const (
	// CompareObjConfRstNoNeed 对比对象配置结果-不需更新
	CompareObjConfRstNoNeed CompareObjConfRst = iota

	// CompareObjConfRstNeedReset 对比对象配置结果-需要重置
	CompareObjConfRstNeedReset

	// CompareObjConfRstNeedReplace 对比对象配置结果-需要替换
	CompareObjConfRstNeedReplace
)

// IObject 定义对象接口类型
type IObject interface {
	io.Closer
}

// IConf 定义配置接口类型
type IConf interface{}

// GetObjConfFunc 定义获取对象配置函数类型
type GetObjConfFunc func() (conf IConf, err error)

// CompareObjConfFunc 定义对比对象配置函数类型
type CompareObjConfFunc func(oldConf, newConf IConf) (rst CompareObjConfRst, err error)

// NewObjFunc 定义新建对象函数类型
type NewObjFunc func(conf IConf) (obj IObject, err error)

// ResetObjFunc 定义重置对象函数类型
type ResetObjFunc func(obj IObject, oldConf, newConf IConf) (err error)

var (
	ErrGetObjConfFuncIsNil     = errors.New("get object conf func is nil")
	ErrCompareObjConfFuncIsNil = errors.New("compare object conf func is nil")
	ErrResetObjFuncIsNil       = errors.New("reset object func is nil")
	ErrNewObjFuncIsNil         = errors.New("new object func is nil")
	ErrContainerClosed         = errors.New("container is closed")
	ErrInvalidConfType         = errors.New("invalid conf type")
	ErrInvalidObjectType       = errors.New("invalid object type")
)

// Container 定义了一个可以在配置发生更新时安全替换或重置封装对象的结构类型,
// 这个类型实现了Updater接口, 因此可以直接注册在Conf类型的updaters上,
// 对Conf监听到的配置修改事件做出响应
type Container struct {
	// 保证实例值的读写安全
	mu sync.RWMutex

	// 获取对象配置函数
	getObjConf GetObjConfFunc

	// 对比对象配置函数
	compareObjConf CompareObjConfFunc

	// 创建对象函数
	newObj NewObjFunc

	// 重置对象函数
	resetObj ResetObjFunc

	// 指向当前对象配置数据
	conf IConf

	// 指向当前对象
	obj IObject

	// 对象锁, 保护对象安全回收
	objmus map[IObject]*sync.RWMutex

	// 关闭状态
	closed bool
}

// NewContainer 根据指定要素创建&初始化Container实例, 并返回实例指针,
// 创建过程中遇到要素不充分时将返回对应错误
func NewContainer(getObjConf GetObjConfFunc, compareObjConf CompareObjConfFunc, newObj NewObjFunc, resetObj ResetObjFunc) (ct *Container, err error) {
	if getObjConf == nil {
		err = ErrGetObjConfFuncIsNil
		return
	}

	if newObj == nil {
		err = ErrNewObjFuncIsNil
		return
	}

	ct = &Container{
		getObjConf:     getObjConf,
		compareObjConf: compareObjConf,
		newObj:         newObj,
		resetObj:       resetObj,
		objmus:         make(map[IObject]*sync.RWMutex, 2),
	}

	if err = ct.init(); err != nil {
		err = fmt.Errorf("init: %w", err)
		return
	}

	return
}

func (ct *Container) init() (err error) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	ct.conf, err = ct.getObjConf()
	if err != nil {
		err = fmt.Errorf("get object conf: %w", err)
		return
	}

	ct.obj, err = ct.newObj(ct.conf)
	if err != nil {
		err = fmt.Errorf("new object: %w", err)
		return
	}

	ct.objmus[ct.obj] = &sync.RWMutex{}

	return
}

// MustGetObj 返回Container指向的当前对象, 如果Container已经关闭将导致panic,
// 返回对象的同时, 会对该对象记录读锁, 读锁未释放前, 对象不能初回收
func (ct *Container) MustGetObj() (obj IObject) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	if ct.closed {
		panic(ErrContainerClosed)
	}
	(ct.objmus[ct.obj]).RLock()
	obj = ct.obj

	return
}

// PutObj 释放指定对象, 实际是释放该对象所关联的一把读锁,
// MustGetObj和PutObj应成对出现, 如果MustGetObj获取到的对象在离开作用域前没有释放,
// 将导致对应的锁资源泄露
func (ct *Container) PutObj(obj IObject) {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	mu, ok := ct.objmus[obj]
	if ok {
		mu.RUnlock()
	}

	return
}

// Update 实现Updater接口, 用于注册在配置更新时回调, 方法将通过回调方式替换或重置Container内置对象,
// 需要注意: 为保证系统稳定运行, 对旧对象的回收是异步的, 当前实现忽略了回收时可能发生的错误
func (ct *Container) Update() (err error) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if ct.closed {
		err = ErrContainerClosed
		return
	}

	if ct.compareObjConf == nil {
		err = ErrCompareObjConfFuncIsNil
		return
	}

	nconf, err := ct.getObjConf()
	if err != nil {
		err = fmt.Errorf("get object conf: %w", err)
		return
	}

	rst, err := ct.compareObjConf(ct.conf, nconf)
	if err != nil {
		err = fmt.Errorf("compare object conf: %w", err)
		return
	}

	switch rst {
	case CompareObjConfRstNeedReset:
		if ct.resetObj == nil {
			err = ErrResetObjFuncIsNil
			return
		}

		if err = ct.resetObj(ct.obj, ct.conf, nconf); err != nil {
			err = fmt.Errorf("reset object: %w", err)
			return
		}

		ct.conf = nconf
		return

	case CompareObjConfRstNeedReplace:
		if ct.newObj == nil {
			err = ErrNewObjFuncIsNil
			return
		}

		nobj, nerr := ct.newObj(nconf)
		if nerr != nil {
			err = fmt.Errorf("new object: %w", nerr)
			return
		}

		go func(obj IObject, objmu *sync.RWMutex) {
			objmu.Lock()
			defer objmu.Unlock()
			obj.Close()
		}(ct.obj, ct.objmus[ct.obj])

		ct.conf = nconf
		ct.obj = nobj
		ct.objmus[ct.obj] = &sync.RWMutex{}

	case CompareObjConfRstNoNeed:
		return
	}

	return
}

// Close 实现io.Closer接口, 用于回收Container及Container内置的对象,
// Container回收时, 内置对象的回收不同于在Update中, 该过程是同步的,
// 因此内置对象被回收前需要等待关联的锁资源释放, 关闭失败的错误也会同步返回,
// 同样会返回错误, 回收后Container的closed标记将置为true, 此时所有在Container上的调用将是非法的
func (ct *Container) Close() (err error) {
	ct.mu.Lock()

	if ct.closed {
		err = ErrContainerClosed
		ct.mu.Unlock()
		return
	}

	ct.closed = true
	ct.mu.Unlock()

	(ct.objmus[ct.obj]).Lock()
	defer (ct.objmus[ct.obj]).Unlock()

	if err = ct.obj.Close(); err != nil {
		err = fmt.Errorf("object close: %w", err)
		return
	}

	return
}
