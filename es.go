package goes

import (
	"bytes"
	"fmt"
	"time"

	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"

	"gitlab.superjq.com/go-tools/logger"
)

type Option struct {
	Host, // 链接地址
	Port, // 端口号
	User, // 用户名
	Pass, // 密码
	Scheme string
	Debug int    // 是否开启 debug 日志
	Log   Logger // 日志实例,默认使用 zap log
}

var (
	esclient *elastic.Client
	esopt    *Option
	esLog    Logger
)

func New(opt *Option) *elastic.Client {
	esopt = opt
	if opt.Log == nil {
		esLog = newLogger()
	} else {
		esLog = opt.Log
	}
	return InitEsClient()

}

func InitEsClient() *elastic.Client {
	esurls := getEsUrl()
	sc := esopt.Scheme

	esoptions := getBaseOptions(esopt.User, esopt.Pass, esurls...)

	if esopt.Debug == 1 {
		esoptions = append(esoptions, elastic.SetInfoLog(esLog))
	}

	if len(sc) > 0 {
		esoptions = append(esoptions, elastic.SetScheme(sc))
		esoptions = append(esoptions, elastic.SetHealthcheck(false))
	}

	es, err := elastic.NewClient(esoptions...)
	if err != nil {
		esLog.Printf("new es client err : %s", err.Error())
		panic(err)
	}

	// 测试链接
	for _, eu := range esurls {
		info, code, err := es.Ping(eu).Do(context.Background())
		if err != nil {
			fmt.Println("es ping err: ", err.Error())
			panic(err)
		}
		logger.Infof("Elasticsearch init code %d and version %s\n", code, info.Version.Number)
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
	// elasticsearch 服务地址，多个服务地址放在切片中
	options = append(options, elastic.SetURL(urls...))
	// 基于http base auth验证机制的账号和密码
	options = append(options, elastic.SetBasicAuth(username, password))
	options = append(options, elastic.SetHealthcheckTimeoutStartup(15*time.Second))
	// 开启Sniff，SDK会定期(默认15分钟一次)嗅探集群中全部节点，将全部节点都加入到连接列表中，
	// 后续新增的节点也会自动加入到可连接列表，但实际生产中我们可能会设置专门的协调节点，所以默认不开启嗅探
	options = append(options, elastic.SetSniff(false))
	options = append(options, elastic.SetErrorLog(esLog))
	options = append(options, elastic.SetGzip(true))
	return options

}
