// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Router_Handler_Strict_WithObject(t *testing.T) {
	type TestReq struct {
		Age  int
		Name string
	}
	type TestRes struct {
		Id   int
		Age  int
		Name string
	}
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test", func(ctx context.Context, req *TestReq) (res *TestRes, err error) {
		return &TestRes{
			Id:   1,
			Age:  req.Age,
			Name: req.Name,
		}, nil
	})
	s.BindHandler("/test/error", func(ctx context.Context, req *TestReq) (res *TestRes, err error) {
		return &TestRes{
			Id:   1,
			Age:  req.Age,
			Name: req.Name,
		}, gerror.New("error")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(client.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)
	})
}

type TestForHandlerWithObjectAndMeta1Req struct {
	g.Meta `path:"/custom-test1" method:"get"`
	Age    int
	Name   string
}

type TestForHandlerWithObjectAndMeta1Res struct {
	Id  int
	Age int
}

type TestForHandlerWithObjectAndMeta2Req struct {
	g.Meta `path:"/custom-test2" method:"get"`
	Age    int
	Name   string
}

type TestForHandlerWithObjectAndMeta2Res struct {
	Id   int
	Name string
}

type ControllerForHandlerWithObjectAndMeta1 struct{}

func (ControllerForHandlerWithObjectAndMeta1) Index(ctx context.Context, req *TestForHandlerWithObjectAndMeta1Req) (res *TestForHandlerWithObjectAndMeta1Res, err error) {
	return &TestForHandlerWithObjectAndMeta1Res{
		Id:  1,
		Age: req.Age,
	}, nil
}

func (ControllerForHandlerWithObjectAndMeta1) Test2(ctx context.Context, req *TestForHandlerWithObjectAndMeta2Req) (res *TestForHandlerWithObjectAndMeta2Res, err error) {
	return &TestForHandlerWithObjectAndMeta2Res{
		Id:   1,
		Name: req.Name,
	}, nil
}

type TestForHandlerWithObjectAndMeta3Req struct {
	g.Meta `path:"/custom-test3" method:"get"`
	Age    int
	Name   string
}

type TestForHandlerWithObjectAndMeta3Res struct {
	Id  int
	Age int
}

type TestForHandlerWithObjectAndMeta4Req struct {
	g.Meta `path:"/custom-test4" method:"get"`
	Age    int
	Name   string
}

type TestForHandlerWithObjectAndMeta4Res struct {
	Id   int
	Name string
}

type ControllerForHandlerWithObjectAndMeta2 struct{}

func (ControllerForHandlerWithObjectAndMeta2) Test3(ctx context.Context, req *TestForHandlerWithObjectAndMeta3Req) (res *TestForHandlerWithObjectAndMeta3Res, err error) {
	return &TestForHandlerWithObjectAndMeta3Res{
		Id:  1,
		Age: req.Age,
	}, nil
}

func (ControllerForHandlerWithObjectAndMeta2) Test4(ctx context.Context, req *TestForHandlerWithObjectAndMeta4Req) (res *TestForHandlerWithObjectAndMeta4Res, err error) {
	return &TestForHandlerWithObjectAndMeta4Res{
		Id:   1,
		Name: req.Name,
	}, nil
}

func Test_Router_Handler_Strict_WithObjectAndMeta(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", new(ControllerForHandlerWithObjectAndMeta1))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), `{"code":65,"message":"Not Found","data":null}`)
		t.Assert(client.GetContent(ctx, "/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/custom-test2?age=18&name=john"), `{"code":65,"message":"Not Found","data":null}`)
	})
}

func Test_Router_Handler_Strict_Group_Bind(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(ControllerForHandlerWithObjectAndMeta1),
			new(ControllerForHandlerWithObjectAndMeta2),
		)
	})
	s.Group("/api/v2", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(ControllerForHandlerWithObjectAndMeta1),
			new(ControllerForHandlerWithObjectAndMeta2),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), `{"code":65,"message":"Not Found","data":null}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":65,"message":"Not Found","data":null}`)

		t.Assert(client.GetContent(ctx, "/api/v1/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)

		t.Assert(client.GetContent(ctx, "/api/v2/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
	})
}

type TestForHandlerWithObjectAndMeta5Req struct {
	g.Meta `path:"/test-meta" method:"get"`
	Code   int `p:"code"`
}

type TestForHandlerWithObjectAndMeta5Res struct {
	g.Meta `200:"ok" 400:"bad request" 401:"unauthorized" 500:"internal server error"`
}

type ControllerForHandlerWithObjectAndMeta3 struct{}

func (ControllerForHandlerWithObjectAndMeta3) Test1(ctx2 context.Context, req *TestForHandlerWithObjectAndMeta5Req) (res *TestForHandlerWithObjectAndMeta5Res, err error) {
	return
}

func TestHandlerForReqResMeta(t *testing.T) {
	type TestNoMetaReq struct{}
	type TestNoMetaRes struct{}

	s := g.Server(guid.S())

	// directly response with metadata without executing handler.
	s.BindMiddleware("/*", func(r *ghttp.Request) {
		data := g.Map{
			"reqMeta": r.GetServeHandler().Handler.Info.ReqMeta,
			"resMeta": r.GetServeHandler().Handler.Info.ResMeta,
		}

		r.Response.WriteExit(data)
	})

	s.BindObject("/", new(ControllerForHandlerWithObjectAndMeta3))

	// no g.Meta defined in "XxxReq" and "XxxRes"
	s.BindHandler("/test-meta-empty", func(ctx context.Context, req TestNoMetaReq) (res TestNoMetaRes, err error) {
		return
	})

	// no "XxxReq" and "XxxRes"
	s.BindHandler("/test-meta-classical-handler", func(r *ghttp.Request) {
		return
	})

	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/test-meta"), `{"reqMeta":{"method":"get","path":"/test-meta"},"resMeta":{"200":"ok","400":"bad request","401":"unauthorized","500":"internal server error"}}`)
		t.Assert(client.GetContent(ctx, "/test-meta-empty"), `{"reqMeta":null,"resMeta":null}`)
		t.Assert(client.GetContent(ctx, "/test-meta-classical-handler"), `{"reqMeta":null,"resMeta":null}`)
	})
}

func Test_Issue1708(t *testing.T) {
	type Test struct {
		Name string `json:"name"`
	}
	type Req struct {
		Page       int      `json:"page"       dc:"分页码"`
		Size       int      `json:"size"       dc:"分页数量"`
		TargetType string   `json:"targetType" v:"required#评论内容类型错误" dc:"评论类型: topic/ask/article/reply"`
		TargetId   uint     `json:"targetId"   v:"required#评论目标ID错误" dc:"对应内容ID"`
		Test       [][]Test `json:"test"`
	}
	type Res struct {
		Page       int      `json:"page"       dc:"分页码"`
		Size       int      `json:"size"       dc:"分页数量"`
		TargetType string   `json:"targetType" v:"required#评论内容类型错误" dc:"评论类型: topic/ask/article/reply"`
		TargetId   uint     `json:"targetId"   v:"required#评论目标ID错误" dc:"对应内容ID"`
		Test       [][]Test `json:"test"`
	}

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test", func(ctx context.Context, req *Req) (res *Res, err error) {
		return &Res{
			Page:       req.Page,
			Size:       req.Size,
			TargetType: req.TargetType,
			TargetId:   req.TargetId,
			Test:       req.Test,
		}, nil
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		content := `
{
    "targetType":"topic",
    "targetId":10785,
    "test":[
        [
            {
                "name":"123"
            }
        ]
    ]
}
`
		t.Assert(
			client.PostContent(ctx, "/test", content),
			`{"code":0,"message":"","data":{"page":0,"size":0,"targetType":"topic","targetId":10785,"test":[[{"name":"123"}]]}}`,
		)
	})
}
