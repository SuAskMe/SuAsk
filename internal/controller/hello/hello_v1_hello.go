package hello

import (
	"context"
	"suask/api/hello/v1"
)

func (c *ControllerV1) Hello(ctx context.Context, req *v1.HelloReq) (res *v1.HelloRes, err error) {
	res = &v1.HelloRes{
		Message: "hello",
	}
	return
}
