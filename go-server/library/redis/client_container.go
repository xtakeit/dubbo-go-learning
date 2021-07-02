package redis

import (
	"errors"
	"fmt"
	"go-server/library/conf"
)

var (
	ErrGetClientConfFuncIsNil = errors.New("get client conf func is nil")
)

type ClientContainer struct {
	*conf.Container
}

type GetClientConfFunc func() (*ClientConf, error)

func NewContainer(getCliCf GetClientConfFunc) (ct *ClientContainer, err error) {
	if getCliCf == nil {
		err = ErrGetClientConfFuncIsNil
		return
	}

	getObjConf := func() (icf conf.IConf, err error) {
		icf, err = getCliCf()
		if err != nil {
			err = fmt.Errorf("get client conf: %w", err)
			return
		}

		return
	}

	ict, err := conf.NewContainer(getObjConf, compareClientConf, newClientObj, nil)
	if err != nil {
		err = fmt.Errorf("new conf container: %w", err)
		return
	}

	ct = &ClientContainer{
		Container: ict,
	}

	return
}

func newClientObj(icf conf.IConf) (iobj conf.IObject, err error) {
	cf, ok := icf.(*ClientConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	iobj, err = NewClient(cf)
	if err != nil {
		err = fmt.Errorf("new client: %w", err)
		return
	}

	return
}

func compareClientConf(iocf, incf conf.IConf) (rst conf.CompareObjConfRst, err error) {
	ocf, ok := iocf.(*ClientConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	ncf, ok := incf.(*ClientConf)
	if !ok {
		err = conf.ErrInvalidConfType
		return
	}

	if *ocf != *ncf {
		rst = conf.CompareObjConfRstNeedReplace
		return
	}

	rst = conf.CompareObjConfRstNoNeed

	return
}

func (ct *ClientContainer) MustGetClient() (cli *Client) {
	obj := ct.MustGetObj()

	cli, ok := obj.(*Client)
	if !ok {
		panic(conf.ErrInvalidObjectType)
	}

	return
}

func (ct *ClientContainer) PutClient(cli *Client) {
	ct.PutObj(cli)
	return
}
