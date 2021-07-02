// Author: Steve Zhang
// Date: 2020/9/21 5:04 下午

package kafka

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
)

type GroupConsumer struct {
	mu               sync.Mutex
	consumerGroup    sarama.ConsumerGroup
	topics           []string
	handleMessage    MessageHandleFunc
	handleConsumeErr ConsumeErrHandleFunc
	clean            func()
}

type ConsumeErrHandleFunc func(error)

// ProducerConf 定义同步生产者配置类型
type GroupConsumerConf struct {
	// Kafka集群节点地址, 多个用英文','号隔开
	Brokers string

	// 分组ID
	GroupID string

	// 消费主题, 多个用英文','号隔开
	Topics string

	// 扩展配置, 需要覆盖sarama默认配置时使用
	Ext *sarama.Config
}

func NewGroupConsumer(cf *GroupConsumerConf) (consumer *GroupConsumer, err error) {
	if cf.Ext == nil {
		cf.Ext = sarama.NewConfig()
		cf.Ext.Version = sarama.V2_6_0_0
	}

	cg, err := sarama.NewConsumerGroup(strings.Split(cf.Brokers, ","), cf.GroupID, cf.Ext)

	if err != nil {
		err = fmt.Errorf("sarama.NewConsumerGroup: %w", err)
		return
	}

	consumer = &GroupConsumer{
		consumerGroup: cg,
		topics:        strings.Split(cf.Topics, ","),
		handleMessage: func(message *sarama.ConsumerMessage) {
			// 默认处理函数不会对消息做任何处理
		},
	}

	return
}

// 需要在Run前调用, fn不能为nil
func (gc *GroupConsumer) SetMessageHandleFunc(f MessageHandleFunc) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.handleMessage = f
}

func (gc *GroupConsumer) SetConsumeErrHandleFunc(f ConsumeErrHandleFunc) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.handleConsumeErr = f
}

// Start 启动分组消费, 调用前应先设置消息处理函数, 否则将返回错误
func (gc *GroupConsumer) Start() {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	handler := NewGroupConsumerHandler(gc.handleMessage)
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := gc.consumerGroup.Consume(ctx, gc.topics, handler); err != nil && gc.handleConsumeErr != nil {
				gc.handleConsumeErr(err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	gc.clean = func() {
		cancel()
		wg.Wait()
	}

	return
}

func (gc *GroupConsumer) Close() (err error) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	if gc.clean == nil {
		return
	}

	gc.clean()

	if err = gc.consumerGroup.Close(); err != nil {
		err = fmt.Errorf("sarama.ConsumerGroup.Close: %w", err)
		return
	}

	return
}
