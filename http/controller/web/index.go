package web

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
)

type Index struct {
}

func (i *Index) Index(c *gin.Context) {
	c.Redirect(302, "/_admin/")
}

func (i *Index) ConfigJs(c *gin.Context) {
	jsString := func(value string) string {
		encoded, _ := json.Marshal(value)
		return string(encoded)
	}

	apiServer := jsString(global.Config.Rustdesk.ApiServer)
	idServer := jsString(global.Config.Rustdesk.IdServer)
	key := jsString(global.Config.Rustdesk.Key)
	wsHost := jsString(global.Config.Rustdesk.WsHost)
	magicQueryonline := global.Config.Rustdesk.WebclientMagicQueryonline
	tmp := fmt.Sprintf(`const rustdeskApiServer = %s;
const rustdeskIdServer = %s;
const rustdeskKey = %s;
const rustdeskWsHost = %s;
const ws2_prefix = 'wc-';

localStorage.setItem('api-server', rustdeskApiServer);
localStorage.setItem(ws2_prefix+'api-server', rustdeskApiServer);
if (rustdeskIdServer) {
  localStorage.setItem('custom-rendezvous-server', rustdeskIdServer);
  localStorage.setItem(ws2_prefix+'custom-rendezvous-server', rustdeskIdServer);
}
if (rustdeskKey) {
  localStorage.setItem('key', rustdeskKey);
  localStorage.setItem(ws2_prefix+'key', rustdeskKey);
}

window.webclient_magic_queryonline = %d;
window.ws_host = rustdeskWsHost;
`, apiServer, idServer, key, wsHost, magicQueryonline)

	c.Header("Content-Type", "application/javascript")
	c.String(200, tmp)
}
