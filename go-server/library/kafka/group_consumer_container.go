// Author: Steve Zhang
// Date: 2020/9/22 2:24 下午

package kafka

import (
	"errors"
	"fmt"

	"go-server/library/conf"
)

var (
	ErrGetGroupConsumerConfFuncIsNil = errors.New("get group consumer conf func is nil")
)

type GroupConsumerContainer struct {
	*conf.Container
}

type GetGroupConsumerConfFunc func() (*GroupConsumerConf, error)

func NewGroupConsumerContainer(getGroupConsumerCf GetGroupConsumerConfFunc) (ct *GroupConsumerContainer, err error) {
	if getGroupConsumerCf == nil {
		err = ErrGetGroupConsumerConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getGroupConsumerCf()
		if err != nil {
			err = fmt.Errorf("get group consumer conf: %w", err)
			return
		}

		return
	}

	ict, err := conf.NewContainer(getObjConf, compareClientConf, newClientObj, nil)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &GroupConsumerContainer{
		Container: ict,
	}

	return
}

func newClientObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*GroupConsumerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewGroupConsumer(cf)
	if err != nil {
		err = fmt.Errorf("new group consumer: %w", err)
		return
	}

	return
}

func compareClientConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ok := iocf.(*GroupConsumerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*GroupConsumerConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	switch {
	case ocf.Brokers != ncf.Brokers,
		ocf.GroupID != ncf.GroupID,
		ocf.Topics != ncf.Topics:
		rst = conf.CompareObjConfRstNeedReplace
		return

	}

	rst = conf.CompareObjConfRstNoNeed

	return
}

func (ct *GroupConsumerContainer) MustGetGroupConsumer() (consumer *GroupConsumer) {
	obj := ct.MustGetObj()

	consumer, ok := obj.(*GroupConsumer)
	if !ok {
		panic(conf.ErrInvalidObjectType)
	}

	return
}

func (ct *GroupConsumerContainer) PutGroupConsumer(consumer *GroupConsumer) {
	ct.PutObj(consumer)
	return
}
