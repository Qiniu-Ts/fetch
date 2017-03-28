package cdn

import (
	"qiniupkg.com/api.v7/auth/qbox"
	"qiniupkg.com/x/rpc.v7"
)

type Client struct {
	rpc.Client
}

func NewClient(ak, sk string) *Client {
	mac := &qbox.Mac{ak, []byte(sk)}

	cli := rpc.Client{qbox.NewClient(mac, nil)}

	return &Client{cli}
}

type ListRet struct {
	Data map[string][]LogEntry `json:"data"`
}

type LogEntry struct {
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Mtime int    `json:"mtime"`
	Url   string `json:"url"`
}

func (c *Client) List(day, domains string) (list map[string][]LogEntry, err error) {
	u := "http://fusion.qiniuapi.com/v2/tune/log/list"

	params := map[string]interface{}{"day": day, "domains": domains}
	ret := ListRet{}
	err = c.CallWithJson(nil, &ret, "POST", u, params)
	list = ret.Data

	return
}
