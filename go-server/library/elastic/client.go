// Author: Steve Zhang
// Date: 2020/9/25 3:27 下午

package elastic

import (
	"fmt"

	"github.com/olivere/elastic"
)

type Client struct {
	*elastic.Client
}

type ClientConf struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewClient(cf *ClientConf) (client *Client, err error) {
	var options []elastic.ClientOptionFunc

	options = append(options, elastic.SetSniff(false))

	if cf.URL != "" {
		options = append(options, elastic.SetURL(cf.URL))
	}
	if cf.Username != "" {
		options = append(options, elastic.SetBasicAuth(cf.Username, cf.Password))
	}

	cli, err := elastic.NewClient(options...)

	if err != nil {
		err = fmt.Errorf("elastice.NewClient: %w", err)
		return
	}

	client = &Client{
		Client: cli,
	}

	return
}

func (cli *Client) Close() (err error) {
	return nil
}
