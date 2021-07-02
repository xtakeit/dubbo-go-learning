// Author: Steve Zhang
// Date: 2020/9/16 6:01 下午

package kafka

import (
	"github.com/Shopify/sarama"
)

func (ct *SyncProducerContainer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	pdr := ct.MustGetProducer()
	defer ct.PutProducer(pdr)

	partition, offset, err = pdr.SendMessage(msg)

	return
}

func (ct *SyncProducerContainer) SendMessages(msgs []*sarama.ProducerMessage) (err error) {
	pdr := ct.MustGetProducer()
	defer ct.PutProducer(pdr)

	err = pdr.SendMessages(msgs)

	return
}
