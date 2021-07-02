// Author: Steve Zhang
// Date: 2020/9/9 6:06 下午

package kafka

import (
	"errors"
	"fmt"

	"go-server/library/conf"
)

var (
	ErrGetProducerConfFuncIsNil = errors.New("get producer conf func is nil")
)

type SyncProducerContainer struct {
	*conf.Container
}

type GetProducerConfFunc func() (*ProducerConf, error)

func NewSyncProducerContainer(getPdrCf GetProducerConfFunc) (ct *SyncProducerContainer, err error) {
	if getPdrCf == nil {
		err = ErrGetProducerConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getPdrCf()
		if err != nil {
			err = fmt.Errorf("get producer conf: %w", err)
			return
		}

		return
	}

	ict, err := conf.NewContainer(getObjConf, compareProducerConf, newSyncProducerObj, nil)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &SyncProducerContainer{
		Container: ict,
	}

	return
}

func newSyncProducerObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*ProducerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewSyncProducer(cf)
	if err != nil {
		err = fmt.Errorf("new sync producer: %w", err)
		return
	}

	return
}

func compareProducerConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ok := iocf.(*ProducerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*ProducerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	if ocf.Brokers != ncf.Brokers {
		rst = conf.CompareObjConfRstNeedReplace
		return
	}

	rst = conf.CompareObjConfRstNoNeed

	return
}

func (ct *SyncProducerContainer) MustGetProducer() (pdr *SyncProducer) {
	obj := ct.MustGetObj()

	pdr, ok := obj.(*SyncProducer)
	if !ok {
		panic(conf.ErrInvalidObjectType)
	}

	return
}

func (ct *SyncProducerContainer) PutProducer(pdr *SyncProducer) {
	ct.PutObj(pdr)

	return
}
