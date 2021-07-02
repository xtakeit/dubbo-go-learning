// Author: Steve Zhang
// Date: 2020/9/21 5:49 下午

package kafka

import "github.com/Shopify/sarama"

type MessageHandleFunc func(message *sarama.ConsumerMessage)

func NewGroupConsumerHandler(handleMessage MessageHandleFunc) *GroupConsumerHandler {
	return &GroupConsumerHandler{
		handleMessage: handleMessage,
	}
}

type GroupConsumerHandler struct {
	handleMessage func(message *sarama.ConsumerMessage)
}

func (cgh *GroupConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *GroupConsumerHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	return nil
}

func (cgh *GroupConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for message := range claim.Messages() {
		if cgh.handleMessage != nil {
			cgh.handleMessage(message)
		}

		sess.MarkMessage(message, "")
	}

	return nil
}
