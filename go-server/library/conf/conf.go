// Author: Steve Zhang
// conf 包封装了一些应用配置相关的基础功能, 包括配置文件加载, 配置内容读取及配置热更新等,
// 为方便兼容apollo系统, 目前仅支持唯一的.env类型的配置文件,
// 配置热更新通过监听配置文件句柄事件回调注册方法来实现, 关注的句柄事件为: 创建, 写入,
// 为保证服务稳定, 删除配置文件不会触发配置重载, 但服务重启时可能会因缺少配置文件而失败
package conf

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"

	"github.com/joho/godotenv"
)

var (
	ErrContentNotLoaded   = errors.New("config content is not loaded") // 配置内容未加载错误
	ErrRepeatedlyWatching = errors.New("repeatedly watching")          // 重复启动监听错误
)

// Conf 配置类型定义, 封装了服务配置功能所需的成员与方法
type Conf struct {
	mutex             sync.RWMutex            // 读写锁, 保证并发时封装成员的读写安全
	path              string                  // 配置文件绝对路径
	items             map[string]string       // 配置内容记录
	updaters          []Updater               // 注册更新者, 在配置重载时需要更新的实例列表
	afterLoaded       func(map[string]string) // 勾子函数, 在配置重载后调用, 可以改变配置值
	beforeUpdateHooks []func()                // 勾子函数, 在配置重载updaters刷新前调用
	afterUpdateHooks  []func()                // 勾子函数, 在配置重载updaters刷新后调用
	herr              func(error)             // 配置监听错误处理回调
	exit              chan struct{}           // 监听退出信道
	loaded            bool                    // 配置内容加载状态
	watching          bool                    // 配置文件监听状态
}

// Updater 更新类型接口, 注册在配置重载时进行更新的实例的接口约束
type Updater interface {
	Update() error
}

// NewConf 返回绑定到指定配置文件envPath的Conf指针cf和错误err,
// 当获取filename绝对路径失败时, 将返回nil的cf和非nil的err, 否则将返回初始化的cf和nil的err
func NewConf(envPath string) (cf *Conf, err error) {
	path, err := filepath.Abs(envPath)
	if err != nil {
		err = fmt.Errorf("abstract filepath: %w", err)
		return
	}

	cf = &Conf{
		path: path,
		exit: make(chan struct{}),
	}

	return
}

// Load 从Conf实例绑定的配置文件中加载配置内容, 当前仅持.env类型的配置文件, 且读取失败时将返回错误,
// 该方法在首次执行前调用其它读取配置的方法如Scan, Get等将返回配置内容未加载的错误,
// 该方法可安全并发且重复调用, 但加载配置期间其它读取配置的协程将阻塞直到写锁释放,
// 可通过重复调用该方法来读取最新的配置文件信息, 重复调用失败时不会影响到已加载的配置内容
func (cf *Conf) Load() (err error) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	items, err := godotenv.Read(cf.path)
	if err != nil {
		err = fmt.Errorf("read env file %s: %w", cf.path, err)
		return
	}

	if cf.afterLoaded != nil {
		cf.afterLoaded(items)
	}

	cf.items = items
	cf.loaded = true

	return
}

// Scan 将配置内容以指定成员标签tag, 扫描到传入的结构体成员指针st上, st的类型必须是非nil的结构体指针,
// 否则将返回错误. 扫描时通过cf.items的key值与扫描对象结构体成员tag的值进行对应,
// 当存在对应tag但在items里没有对应的配置项时, 将会返回错误,
// 当实例未调用Load或Load失败时直接调用Scan将返回内容未加载错误
func (cf *Conf) Scan(st interface{}, tag string) (err error) {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if !cf.loaded {
		err = ErrContentNotLoaded
		return
	}

	if err = scan(cf.items, st, tag); err != nil {
		err = fmt.Errorf("scan map: %w", err)
		return
	}

	return
}

// Get从Conf实例的items中返回指定key的值val, key是否存在ok及错误err,
// 实例未调用Load或Load失败时直接调用Get将返回内容未加载错误
func (cf *Conf) Get(key string) (val string, ok bool, err error) {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if !cf.loaded {
		err = ErrContentNotLoaded
		return
	}

	val, ok = cf.items[key]

	return
}

// MustGet 从配置中获取指定key的值, 如果值不存在或为空将导致panic
func (cf *Conf) MustGet(key string) (val string) {
	cf.mutex.RLock()
	defer cf.mutex.RUnlock()

	if !cf.loaded {
		panic("config unloaded")
		return
	}

	val, ok := cf.items[key]
	if !ok || val == "" {
		panic("config " + key + " undefined")
		return
	}

	return
}

// MustGetBool 获取int类型的配置值
func (cf *Conf) MustGetInt(key string) (val int) {
	valStr := cf.MustGet(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		panic("config " + key + " convert to int failed")
	}
	return
}

// MustGetBool 获取bool类型的配置值
func (cf *Conf) MustGetBool(key string) (val bool) {
	valStr := cf.MustGet(key)
	switch strings.ToLower(valStr) {
	case "false", "", "0":
		val = false
	default:
		val = true
	}
	return
}

// PushUpdater 向Conf实例上注册更新者实例updater, updater实现了Updater接口,
// 将在配置重载时调用Update方法, 当Conf实例上注册了监听错误处理函数, Update方法返回的错误将由该函数处理,
// 当有多个Updater注册到Conf实例上时, updater的调用顺序和注册时的顺序相反
func (cf *Conf) PushUpdater(updater Updater) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.updaters = append(
		cf.updaters, updater,
	)
}

func (cf *Conf) RegisterAfterLoadedHook(hook func(map[string]string)) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.afterLoaded = hook
}

func (cf *Conf) RegisterBeforeUpdateHook(hook func()) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.beforeUpdateHooks = append(cf.beforeUpdateHooks, hook)
}

func (cf *Conf) RegisterAfterUpdateHook(hook func()) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.afterUpdateHooks = append(cf.afterUpdateHooks, hook)
}

// SetWatchErrHandleFunc 为Conf实例注册监听错误处理函数, 用于处理监听过程中发生的所有错误
func (cf *Conf) SetWatchErrHandleFunc(f func(error)) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	cf.herr = f
}

// Watch 创建一个文件句柄Watcher, 并监听Conf实例绑定的配置文件path的目录,
// 当配置文件发生对应的创建和写入事件时, 重载配置内容并按先入后出的顺序执行注册在Conf实例上的refresher的Refresh方法,
// 每个Conf实例同一时间只能启动一个Watch, 启动后以通过CloseWatch方法退出Watch, 重复启动Watch将返回ErrRepeatedlyWatching错误,
// watcher监听的是配置文件的目录, 所以当删除一个配置文件后重新创建这个配置文件也会触发配置的重载与刷新
func (cf *Conf) Watch() (err error) {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	if cf.watching {
		err = ErrRepeatedlyWatching
		return
	}

	wc, err := fsnotify.NewWatcher()
	if err != nil {
		err = fmt.Errorf("new watcher: %w", err)
		return
	}

	dir := filepath.Dir(cf.path)
	if err = wc.Add(dir); err != nil {
		err = fmt.Errorf("watch add %s: %w", dir, err)
		return
	}

	go watch(cf, wc)
	cf.watching = true

	return
}

// CloseWatch 关闭正在进行的Watch, 如果Watch未启动则直接返回
func (cf *Conf) CloseWatch() {
	cf.mutex.Lock()
	defer cf.mutex.Unlock()

	if !cf.watching {
		return
	}

	cf.exit <- struct{}{}
	cf.watching = false

	return
}

// Close 实现io.Closer接口, 当前仅尝试关闭开启的Watch
func (cf *Conf) Close() (err error) {
	cf.CloseWatch()
	return
}

// watch 使用句柄监听器wc监听Conf实例cf的配置文件句柄事件, 重载cf配置内容并回调注册在cf上的Updater的Update方法
func watch(cf *Conf, wc *fsnotify.Watcher) {
	base := filepath.Base(cf.path)

	for {
		select {
		case ev := <-wc.Events:
			envBase := filepath.Base(ev.Name)
			if envBase != base {
				continue
			}

			if ev.Op&fsnotify.Write != fsnotify.Write && ev.Op&fsnotify.Create != fsnotify.Create {
				continue
			}
			if err := cf.Load(); err != nil && cf.herr != nil {
				err = fmt.Errorf("load config: %w", err)
				cf.herr(err)
				break
			}

			func() {
				cf.mutex.RLock()
				defer cf.mutex.RUnlock()

				for _, hook := range cf.beforeUpdateHooks {
					hook()
				}

				for _, updater := range cf.updaters {
					if err := updater.Update(); err != nil && cf.herr != nil {
						err = fmt.Errorf("update: %w", err)
						cf.herr(err)
					}
				}

				for _, hook := range cf.afterUpdateHooks {
					hook()
				}

			}()
		case err := <-wc.Errors:
			if cf.herr != nil {
				err = fmt.Errorf("watch file: %w", err)
				cf.herr(err)
			}
		case <-cf.exit:
			if err := wc.Close(); err != nil && cf.herr != nil {
				err = fmt.Errorf("close watcher: %w", err)
				cf.herr(err)
			}
			return
		}
	}
}
