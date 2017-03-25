# fetch

## 下载地址
[Linux](http://devtools.qiniu.com/linux64/fetch)
[darwin](http://devtools.qiniu.com/darwin/fetch)


## 运行命令

```
fetch <path_to_config>
```

## 配置文件格式

```
{
	"access_key": "<your_ak>",
	"secret_key": "<your_sk>",
	"bucket_to": "logs",
	"domains": "www.qiniu.com",
	"from": "-w1",
	"to": ""
}
```

其中from，to 表示时间

* 格式为 "2016-01-02"
* "-w<weeks_from_now>" 比如 "-w2" 表示从现在开始向前推两周
* "-w<days_from_now>" 比如 "-d4" 表示从现在开始向前推四天
