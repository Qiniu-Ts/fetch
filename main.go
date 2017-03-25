package main


import(
	"log"
	"encoding/json"
	"os"
	"cdn"
	"qiniupkg.com/api.v7/kodo"
	"time"
	"strings"
	"strconv"
)


/*
{
	"access_key": "<access_key>",
	"secret_key": "<secret_key>",
	"bucket_to": "<bucket_name>",
	"domains": "<domains>",
	"from": "",  // date 2017-01-11  or -d<days_ago> or -w<weeks_ago>
	"to": ""  // date 2017-03-20   or -d<days_ago> or -w<weeks_ago> or null express today
}
 */

func main() {

	if len(os.Args) < 2 {
		log.Println("fetch <config_file_path>")
		return
	}

	fp := os.Args[1]
	cfg := Config{}

	err := LoadConfig(fp, &cfg)
	if err != nil {
		log.Println(err)
	}
	//log.Println(cfg.From)

	c := NewClient(cfg)
	c.SyncLogs()
}

type Client struct {
	*kodo.Client
	Config
}

func NewClient(cfg Config) *Client  {
	cli :=  &Client{Config: cfg}

	kodo.SetMac(cfg.AccessKey, cfg.SecretKey)
	cli.Client = kodo.NewWithoutZone(nil)
	return  cli
}

func (c *Client) SyncLogs()  {

	logs, err := c.ListLogs()
	if err != nil {
		log.Println("ListLogs Failed", err)
		return
	}

	log.Println("Start...")

	for _, i := range logs {

		if !c.needFetch(i) {
			continue
		}

		err = c.Bucket(c.BucketTo).Fetch(nil, i.Name, i.Url)
		if err != nil {
			log.Println("Fetch Failed", err)
			continue
		}

		log.Println("Fetch Success", i.Name)
	}

	log.Println("Done...")
}

func (c *Client)ListLogs() (entries []cdn.LogEntry, err error) {

	cdnCli := cdn.NewClient(c.AccessKey, c.SecretKey)


	for _, day := range c.Days() {

		log.Println(day)
		var domainEntries map[string][]cdn.LogEntry
		domainEntries, err = cdnCli.List(day, c.Domains)
		if err != nil {
			return
		}

		for _, log := range domainEntries {
			entries = append(entries, log...)
		}
	}

	return
}

func (c *Client)Days() []string  {

	tf := dateParse(c.From)
	tt := time.Now()
	if c.To != "" {
		tt = dateParse(c.To)
	}

	days := []string{}
	for tt.Sub(tf).Hours() >= 0 {
		days = append(days, tf.Format("2006-01-02"))
		tf = tf.AddDate(0, 0, 1)
	}

	return days
}

func dateParse(d string) time.Time  {

	if strings.HasPrefix(d, "-d") {
		days, _ := strconv.Atoi("-" + d[2:])
		return time.Now().AddDate(0,0, days)
	}

	if strings.HasPrefix(d, "-w") {
		weeks, _ := strconv.Atoi("-" + d[2:])
		return time.Now().AddDate(0,0, weeks*7)
	}

	t, _ := time.Parse("2006-01-02", d)

	return  t
}
func (c *Client) needFetch(le cdn.LogEntry) bool {

	b := c.Bucket(c.BucketTo)
	stat, _ := b.Stat(nil, le.Name)
	return  le.Size != stat.Fsize

}

type Config struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	BucketTo string `json:"bucket_to"`
	Domains string `json:"domains"`
	From string `json:"from"`
	To string `json:"to"`
}

func LoadConfig(fname string, cfg *Config) error {
	fh, err := os.Open(fname)
	if err != nil {
		return err
	}

	jparser := json.NewDecoder(fh)
	return  jparser.Decode(&cfg)
}
