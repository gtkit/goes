# goes

- elasticsearch 实例化，boolquery 查询

- 初始化方法

```
import (
	"github.com/gtkit/goes"
	"github.com/gtkit/logger"
	"github.com/olivere/elastic/v7"

	"package/config"
)

func initEsClient() *elastic.Client {
	return goes.New(&goes.Option{
		Host:   config.GetString("es.host"),
		Port:   config.GetString("es.port"),
		User:   config.GetString("es.user"),
		Pass:   config.GetString("es.pass"),
		Scheme: config.GetString("es.scheme"),
		Debug:  config.GetInt("elasticsearch.debug"),
		Log:    logger.EsLogger(), // 想要输出日志,自己传入, 在此使用 zap log
	})

}

```

## 内部仓库需要执行如下配置(已废弃)
```
go env -w GOPRIVATE=gitlab.superjq.com

go env -w GOINSECURE=gitlab.superjq.com

git config --global url."http://gitlab.superjq.com:".insteadOf "https://gitlab.superjq.com"
```
