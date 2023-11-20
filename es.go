package goes

import (
	"bytes"
	"net"
	"net/http"
	"time"

	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
)

type Option struct {
	Host, // 链接地址
	Port, // 端口号
	User, // 用户名
	Pass, // 密码
	Scheme string
	Debug int            // 是否开启 debug 日志
	Log   elastic.Logger // 日志实例,默认使用 zap log
}

var (
	esclient *elastic.Client
	esopt    *Option
)

func New(opt *Option) *elastic.Client {
	esopt = opt
	if opt.Host == "" {
		panic("Elasticsearch host is empty")
	}

	if opt.Port == "" {
		panic("Elasticsearch port is empty")
	}

	return InitEsClient()

}

func InitEsClient() *elastic.Client {
	esurls := getEsUrl()

	esoptions := getBaseOptions(esopt.User, esopt.Pass, esurls...)

	if esopt.Debug == 1 && esopt.Log != nil {
		esoptions = append(esoptions, elastic.SetInfoLog(esopt.Log))
	}

	if len(esopt.Scheme) > 0 {
		esoptions = append(esoptions, elastic.SetScheme(esopt.Scheme))
		esoptions = append(esoptions, elastic.SetHealthcheck(false))
	}

	es, err := elastic.NewClient(esoptions...)

	if err != nil {
		log("New Elasticsearch client err : %s", err.Error())
		panic(err)
	}

	// 测试链接
	for _, eu := range esurls {
		info, code, err := es.Ping(eu).Do(context.Background())
		if err != nil {
			log("Elasticsearch ping err: %s", err.Error())
			panic(err)
		}
		log("Elasticsearch init Success code %d and version %s\n", code, info.Version.Number)
	}
	esclient = es

	// esversionCode, err := Es.ElasticsearchVersion(esurl)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("es version %s\n", esversionCode)
	return es
}

func getEsUrl() []string {
	EsHost := esopt.Host
	EsPort := esopt.Port

	var conn bytes.Buffer // bytes.buffer是一个缓冲byte类型的缓冲器存放着都是byte
	conn.WriteString(EsHost)
	conn.WriteString(":")
	conn.WriteString(EsPort)
	return []string{conn.String()}
}

func getBaseOptions(username, password string, urls ...string) []elastic.ClientOptionFunc {
	options := make([]elastic.ClientOptionFunc, 0)
	httpClient := &http.Client{}
	httpClient.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second, // 连接超时时间
			KeepAlive: 30 * time.Second, // 长连接超时时间
		}).DialContext,

		MaxIdleConnsPerHost:   100,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	options = append(options, elastic.SetHttpClient(&http.Client{}))
	// elasticsearch 服务地址，多个服务地址放在切片中
	options = append(options, elastic.SetURL(urls...))
	// 基于http base auth验证机制的账号和密码
	options = append(options, elastic.SetBasicAuth(username, password))
	options = append(options, elastic.SetHealthcheckTimeoutStartup(15*time.Second))
	// 开启Sniff，SDK会定期(默认15分钟一次)嗅探集群中全部节点，将全部节点都加入到连接列表中，
	// 后续新增的节点也会自动加入到可连接列表，但实际生产中我们可能会设置专门的协调节点，所以默认不开启嗅探
	options = append(options, elastic.SetSniff(false))
	if esopt.Log != nil {
		options = append(options, elastic.SetErrorLog(esopt.Log))
	}
	options = append(options, elastic.SetGzip(true))
	return options

}

func log(format string, v ...interface{}) {
	if esopt.Log != nil {
		esopt.Log.Printf(format, v)
	}
}
