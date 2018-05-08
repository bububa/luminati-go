package luminati

import (
	"encoding/base64"
	"fmt"
	"github.com/levigross/grequests"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

const (
	SUPER_PROXY            = "zproxy.lum-superproxy.io"
	DEFAULT_PORT           = 22225
	DEFAULT_FAILURES_LIMIT = 10
	DEFAULT_REQUESTS_LIMIT = 10
)

var (
	UserAgents = []string{
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
		"Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
		"Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
		"Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
		"Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
		"Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
		"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52",
	}
	REQUEST_TIMEOUT = 30 * time.Second
)

func randomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return UserAgents[rand.Intn(len(UserAgents))]
}

type Client struct {
	username      string
	password      string
	superProxy    string
	port          uint
	requestsLimit uint
	failuresLimit uint
	requests      uint
	failures      uint
	session       *grequests.Session
}

func NewClient(username string, password string) *Client {
	//superProxy, err := SuperProxy()
	//if err != nil {
	//	return nil, err
	//}
	client := &Client{
		username:      username,
		password:      password,
		superProxy:    SUPER_PROXY,
		port:          DEFAULT_PORT,
		requestsLimit: DEFAULT_REQUESTS_LIMIT,
		failuresLimit: DEFAULT_FAILURES_LIMIT,
	}
	client.NewSession()
	return client
}

func (c *Client) SetFailuresLimit(limit uint) {
	c.failuresLimit = limit
}

func (c *Client) Get(target string, ro *grequests.RequestOptions) (resp *grequests.Response, err error) {
	for {
		if c.requests >= c.requestsLimit && c.requestsLimit > 0 {
			c.NewSession()
		}
		c.requests += 1
		resp, err = c.session.Get(target, ro)
		if err != nil {
			c.failures += 1
			c.NewSession()
		} else if !resp.Ok {
			c.failures += 1
			c.NewSession()
		} else {
			return resp, nil
		}
		if c.failures >= c.failuresLimit && c.failuresLimit > 0 {
			c.Reset()
			return nil, err
		}
	}
	return
}

func (c *Client) proxyUrl(sessionId string) string {
	return fmt.Sprintf("http://%s-country-jp-session-%s:%s@%s:%d", c.username, sessionId, c.password, c.superProxy, c.port)
}

func (c *Client) basicAuth(sessionId string) string {
	auth := fmt.Sprintf("%s-country-jp-session-%s:%s", c.username, sessionId, c.password)
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}

func (c *Client) NewSession() {
	c.requests = 0
	rand.Seed(time.Now().UnixNano())
	s := rand.Intn(10000-100) + 100
	sessionId := strconv.Itoa(s)
	proxyURL, _ := url.Parse(c.proxyUrl(sessionId))
	ro := &grequests.RequestOptions{
		Proxies:   map[string]*url.URL{proxyURL.Scheme: proxyURL},
		UserAgent: randomUserAgent(),
		Headers: map[string]string{
			"Proxy-Authorization": c.basicAuth(sessionId),
			"Proxy-Connection":    "Keep-Alive",
		},
	}
	c.session = grequests.NewSession(ro)

}

func (c *Client) Reset() {
	c.NewSession()
	c.failures = 0
}
