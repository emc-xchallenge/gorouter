package router

import (
	. "launchpad.net/gocheck"
	"net/http"
)

type RegistrySuite struct {
	*Registry
}

var _ = Suite(&RegistrySuite{})

var fooReg = &registerMessage{
	Host: "192.168.1.1",
	Port: 1234,
	Uris: []Uri{"foo.vcap.me", "fooo.vcap.me"},
	Tags: map[string]string{
		"runtime":   "ruby18",
		"framework": "sinatra",
	},
	Dea: "",
	App: "12345",
}

var barReg = &registerMessage{
	Host: "192.168.1.2",
	Port: 4321,
	Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	Tags: map[string]string{
		"runtime":   "javascript",
		"framework": "node",
	},
	Dea: "",
	App: "54321",
}

var bar2Reg = &registerMessage{
	Host: "192.168.1.3",
	Port: 1234,
	Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	Tags: map[string]string{
		"runtime":   "javascript",
		"framework": "node",
	},
	Dea: "",
	App: "54321",
}

var emptyReg = &registerMessage{}

func (s *RegistrySuite) SetUpTest(c *C) {
	s.Registry = NewRegistry()
}

func (s *RegistrySuite) TestRegister(c *C) {
	s.Register(fooReg)
	c.Check(s.NumUris(), Equals, 2)

	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 4)

	s.Register(emptyReg)
	c.Check(s.NumUris(), Equals, 4)
}

func (s *RegistrySuite) TestRegisterIgnoreDuplicates(c *C) {
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumInstances(), Equals, 1)

	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumInstances(), Equals, 1)

	s.Unregister(barReg)
	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumInstances(), Equals, 0)
}

func (s *RegistrySuite) TestRegisterUppercase(c *C) {
	m1 := &registerMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	m2 := &registerMessage{
		Host: "192.168.1.1",
		Port: 1235,
		Uris: []Uri{"FOO.VCAP.ME"},
	}

	s.Register(m1)
	s.Register(m2)

	c.Check(s.NumUris(), Equals, 1)
}

func (s *RegistrySuite) TestUnregister(c *C) {
	s.Register(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumInstances(), Equals, 1)

	s.Register(bar2Reg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumInstances(), Equals, 2)

	s.Unregister(barReg)
	c.Check(s.NumUris(), Equals, 2)
	c.Check(s.NumInstances(), Equals, 1)

	s.Unregister(bar2Reg)
	c.Check(s.NumUris(), Equals, 0)
	c.Check(s.NumInstances(), Equals, 0)
}

func (s *RegistrySuite) TestUnregisterUppercase(c *C) {
	m1 := &registerMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	m2 := &registerMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"FOO.VCAP.ME"},
	}

	s.Register(m1)
	s.Unregister(m2)

	c.Check(s.NumUris(), Equals, 0)
}

func (s *RegistrySuite) TestLookup(c *C) {
	m := &registerMessage{
		Host: "192.168.1.1",
		Port: 1234,
		Uris: []Uri{"foo.vcap.me"},
	}

	s.Register(m)

	m1 := s.Lookup(&http.Request{Host: "foo.vcap.me"})
	c.Check(len(m1), Equals, 1)
	c.Check(m1[0], Equals, m.InstanceId())

	m2 := s.Lookup(&http.Request{Host: "FOO.VCAP.ME"})
	c.Check(len(m2), Equals, 1)
	c.Check(m2[0], Equals, m.InstanceId())
}

func (s *RegistrySuite) TestLookupDoubleRegister(c *C) {
	m1 := &registerMessage{
		Host: "192.168.1.2",
		Port: 1234,
		Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	}

	m2 := &registerMessage{
		Host: "192.168.1.2",
		Port: 1235,
		Uris: []Uri{"bar.vcap.me", "barr.vcap.me"},
	}

	s.Register(m1)
	s.Register(m2)

	ms := s.Lookup(&http.Request{Host: "bar.vcap.me"})
	c.Check(len(ms), Equals, 2)
}