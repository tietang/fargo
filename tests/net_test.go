package fargo_test

// MIT Licensed (see README.md) - Copyright (c) 2013 Hudl <@Hudl>

import (
	"github.com/hudl/fargo"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetApps(t *testing.T) {
	e, _ := fargo.NewConnFromConfigFile("./config_sample/local.gcfg")
	Convey("Pull applications", t, func() {
		a, _ := e.GetApps()
		So(len(a["EUREKA"].Instances), ShouldEqual, 2)
	})
	Convey("Pull single application", t, func() {
		a, _ := e.GetApp("EUREKA")
		So(len(a.Instances), ShouldEqual, 2)
		for idx, ins := range a.Instances {
			if ins.HostName == "node1.localdomain" {
				So(a.Instances[idx].IPAddr, ShouldEqual, "172.16.0.11")
				So(a.Instances[idx].HostName, ShouldEqual, "node1.localdomain")
			} else {
				So(a.Instances[idx].IPAddr, ShouldEqual, "172.16.0.22")
				So(a.Instances[idx].HostName, ShouldEqual, "node2.localdomain")
			}
		}
	})
}

func TestRegistration(t *testing.T) {
	e, _ := fargo.NewConnFromConfigFile("./config_sample/local.gcfg")
	i := fargo.Instance{
		HostName:         "i-123456",
		Port:             9090,
		App:              "TESTAPP",
		IPAddr:           "127.0.0.10",
		VipAddress:       "127.0.0.10",
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		SecureVipAddress: "127.0.0.10",
		Status:           fargo.UP,
	}
	Convey("Fail to heartbeat a non-existent instance", t, func() {
		j := fargo.Instance{
			HostName:         "i-6543",
			Port:             9090,
			App:              "TESTAPP",
			IPAddr:           "127.0.0.10",
			VipAddress:       "127.0.0.10",
			DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
			SecureVipAddress: "127.0.0.10",
			Status:           fargo.UP,
		}
		err := e.HeartBeatInstance(&j)
		So(err, ShouldNotBeNil)
	})
	Convey("Register an instance to TESTAPP", t, func() {
		Convey("Instance registers correctly", func() {
			err := e.RegisterInstance(&i)
			So(err, ShouldBeNil)
		})
		Convey("Instance can check in", func() {
			err := e.HeartBeatInstance(&i)
			So(err, ShouldBeNil)
		})
	})
}

func TestConnectionCreation(t *testing.T) {
	Convey("Pull applications", t, func() {
		cfg, err := fargo.ReadConfig("./config_sample/local.gcfg")
		So(err, ShouldBeNil)
		e := fargo.NewConnFromConfig(cfg)
		apps, err := e.GetApps()
		So(err, ShouldBeNil)
		So(len(apps["EUREKA"].Instances), ShouldEqual, 2)
	})
}

func TestMetadataReading(t *testing.T) {
	cfg, err := fargo.ReadConfig("./config_sample/local.gcfg")
	So(err, ShouldBeNil)
	e := fargo.NewConnFromConfig(cfg)
	Convey("Read empty instance metadata", t, func() {
		a, err := e.GetApp("EUREKA")
		So(err, ShouldBeNil)
		i := a.Instances[0]
		_, err = i.Metadata.GetString("SomeProp")
		So(err, ShouldNotBeNil)
	})
	Convey("Read valid instance metadata", t, func() {
		a, err := e.GetApp("TESTAPP")
		So(err, ShouldBeNil)
		So(len(a.Instances), ShouldBeGreaterThan, 0)
		if len(a.Instances) == 0 {
			return
		}
		i := a.Instances[0]
		err = e.AddMetadataString(i, "SomeProp", "AValue")
		So(err, ShouldBeNil)
		v, err := i.Metadata.GetString("SomeProp")
		So(err, ShouldBeNil)
		So(v, ShouldEqual, "AValue")
	})
}
