package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-mq/api/prometheus"
	"go-mq/pool"
	"go-mq/proto"
	"log"
	"net/http"
)

type HttpServer struct {
	Port int
}

func (this *HttpServer) Run() {
	engine := gin.Default()
	engine.Use(
		prometheus.NewBuilder(
			"gateway", "mq_monitor", "gin_http", "", "统计 gin 的http接口数据",
		).BuildResponseTime(),
		prometheus.NewBuilder(
			"gateway", "mq_monitor", "gin_http", "", "统计 gin 的http接口数据",
		).BuildActiveRequest(),
	)
	engine.POST("/apply_mq", this.ApplyMq)
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	err := engine.Run(fmt.Sprintf(":%d", this.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func (this *HttpServer) unmarshal(ctx *gin.Context, result interface{}) {
	buff, err := ctx.GetRawData()
	if err != nil {
		fmt.Println("unmarshal err: ", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)

	} else if err = json.Unmarshal(buff, result); err != nil {
		fmt.Println("unmarshal err: ", err.Error())
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
}

func (this *HttpServer) ApplyMq(ctx *gin.Context) {
	req := &proto.ApplyMqReq{}
	this.unmarshal(ctx, req)
	if len(req.Ips) <= 0 {
		ctx.JSON(http.StatusOK, &proto.CommonResp{
			Code: 404,
			Msg:  "len < 0",
		})
		return
	}
	resource, err := pool.PoolInstance.Get(req.Ips)
	if err != nil {
		fmt.Println(fmt.Sprintf("ApplyMq error: %s, Req:%+v", err.Error(), req))
		ctx.String(http.StatusOK, "")
		return
	}

	resp := &proto.ApplyMqResp{
		SendPort: resource.SendPort,
		RevPort:  resource.RevPort,
	}
	ctx.JSON(http.StatusOK, &proto.CommonResp{
		Data: resp,
	})
}
