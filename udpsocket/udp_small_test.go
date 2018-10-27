// +build small

package UDPSocket

import (
	"bytes"
	"encoding/gob"
	"errors"
	"testing"
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"

	. "github.com/smartystreets/goconvey/convey"
)

var mockMts = []plugin.MetricType{
	*plugin.NewMetricType(core.NewNamespace("foo"), time.Now(), nil, "", 99),
}

func TestMetaData(t *testing.T) {
	Convey("Meta returns proper metadata", t, func() {
		meta := Meta()
		So(meta, ShouldNotBeNil)
		So(meta.Name, ShouldResemble, PluginName)
		So(meta.Version, ShouldResemble, PluginVersion)
		So(meta.Type, ShouldResemble, PluginType)
	})
}

func TestUDPPublisher(t *testing.T) {
	Convey("Create a UDPSocket Publisher", t, func() {
		udpp := NewUdpPublisher()
		Convey("so udp socket publisher should not be nil", func() {
			So(udpp, ShouldNotBeNil)
		})
		Convey("so udp socket publisher should be of publisher plugin type", func() {
			So(udpp, ShouldHaveSameTypeAs, &udpPublisher{})
		})

		configPolicy, err := udpp.GetConfigPolicy()

		Convey("Test GetConfigPolicy()", func() {
			Convey("So config policy should not be nil", func() {
				So(configPolicy, ShouldNotBeNil)
			})
			Convey("So getting a config policy should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("So config policy should be a cpolicy.ConfigPolicy type", func() {
				So(configPolicy, ShouldHaveSameTypeAs, &cpolicy.ConfigPolicy{})
			})
		})
		Convey("Publish content to file", func() {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(mockMts)

			config := make(map[string]ctypes.ConfigValue)
			config["IP_Port"] = ctypes.ConfigValueStr{Value: "127.0.0.1:1000"}

			Convey("invalid contentType", func() {
				err := udpp.Publish("", buf.Bytes(), config)
				So(err, ShouldResemble, errors.New("Unknown content type ''"))
			})
			Convey("empty content", func() {
				err = udpp.Publish(plugin.SnapGOBContentType, []byte{}, config)
				So(err, ShouldNotBeNil)
			})
			Convey("successful publishing", func() {
				err = udpp.Publish(plugin.SnapGOBContentType, buf.Bytes(), config)
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestFormatMetricTypes(t *testing.T) {
	Convey("FormatMetricTypes returns metrics to publish", t, func() {
		metrics := formatMetricTypes(mockMts)
		So(metrics, ShouldNotBeEmpty)
		// formatted metric has namespace represented as a single string
		So(metrics[0].Namespace, ShouldEqual, mockMts[0].Namespace().String())
	})
}
