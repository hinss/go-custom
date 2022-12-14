package redis

import (
	"context"
	"github.com/hinss/go-custom/framework/provider/config"
	"github.com/hinss/go-custom/framework/provider/log"
	tests "github.com/hinss/go-custom/test"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestCustomService_Load(t *testing.T) {
	container := tests.InitBaseContainer()
	container.Bind(&config.CustomConfigProvider{})
	container.Bind(&log.CustomLogServiceProvider{})

	Convey("test get client", t, func() {
		CustomRedis, err := NewCustomRedis(container)
		So(err, ShouldBeNil)
		service, ok := CustomRedis.(*CustomRedis)
		So(ok, ShouldBeTrue)
		client, err := service.GetClient(WithConfigPath("redis.write"))
		So(err, ShouldBeNil)
		So(client, ShouldNotBeNil)
		ctx := context.Background()
		err = client.Set(ctx, "foo", "bar", 1*time.Hour).Err()
		So(err, ShouldBeNil)
		val, err := client.Get(ctx, "foo").Result()
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "bar")
		err = client.Del(ctx, "foo").Err()
		So(err, ShouldBeNil)
	})
}
