// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package builtin_test

import (
	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/interfaces"
	"github.com/snapcore/snapd/interfaces/apparmor"
	"github.com/snapcore/snapd/interfaces/builtin"
	"github.com/snapcore/snapd/interfaces/udev"
	"github.com/snapcore/snapd/snap"
	"github.com/snapcore/snapd/testutil"
)

type HardwareRandomObserveInterfaceSuite struct {
	iface    interfaces.Interface
	slotInfo *snap.SlotInfo
	slot     *interfaces.ConnectedSlot
	plugInfo *snap.PlugInfo
	plug     *interfaces.ConnectedPlug
}

var _ = Suite(&HardwareRandomObserveInterfaceSuite{
	iface: builtin.MustInterface("hardware-random-observe"),
})

const hardwareRandomObserveConsumerYaml = `name: consumer
apps:
 app:
  plugs: [hardware-random-observe]
`

const hardwareRandomObserveCoreYaml = `name: core
type: os
slots:
  hardware-random-observe:
`

func (s *HardwareRandomObserveInterfaceSuite) SetUpTest(c *C) {
	s.plug, s.plugInfo = MockConnectedPlug(c, hardwareRandomObserveConsumerYaml, nil, "hardware-random-observe")
	s.slot, s.slotInfo = MockConnectedSlot(c, hardwareRandomObserveCoreYaml, nil, "hardware-random-observe")
}

func (s *HardwareRandomObserveInterfaceSuite) TestName(c *C) {
	c.Assert(s.iface.Name(), Equals, "hardware-random-observe")
}

func (s *HardwareRandomObserveInterfaceSuite) TestSanitizeSlot(c *C) {
	c.Assert(interfaces.BeforePrepareSlot(s.iface, s.slotInfo), IsNil)
	slot := &snap.SlotInfo{
		Snap:      &snap.Info{SuggestedName: "some-snap"},
		Name:      "hardware-random-observe",
		Interface: "hardware-random-observe",
	}
	c.Assert(interfaces.BeforePrepareSlot(s.iface, slot), ErrorMatches,
		"hardware-random-observe slots are reserved for the core snap")
}

func (s *HardwareRandomObserveInterfaceSuite) TestSanitizePlug(c *C) {
	c.Assert(interfaces.BeforePreparePlug(s.iface, s.plugInfo), IsNil)
}

func (s *HardwareRandomObserveInterfaceSuite) TestAppArmorSpec(c *C) {
	spec := &apparmor.Specification{}
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, s.slot), IsNil)
	c.Assert(spec.SecurityTags(), DeepEquals, []string{"snap.consumer.app"})
	c.Assert(spec.SnippetForTag("snap.consumer.app"), testutil.Contains, "hw_random/rng_{available,current} r,")
}

func (s *HardwareRandomObserveInterfaceSuite) TestUDevSpec(c *C) {
	spec := &udev.Specification{}
	c.Assert(spec.AddConnectedPlug(s.iface, s.plug, s.slot), IsNil)
	c.Assert(spec.Snippets(), HasLen, 2)
	c.Assert(spec.Snippets(), testutil.Contains, `# hardware-random-observe
KERNEL=="hwrng", TAG+="snap_consumer_app"`)
	c.Assert(spec.Snippets(), testutil.Contains, `TAG=="snap_consumer_app", RUN+="/lib/udev/snappy-app-dev $env{ACTION} snap_consumer_app $devpath $major:$minor"`)
}

func (s *HardwareRandomObserveInterfaceSuite) TestStaticInfo(c *C) {
	si := interfaces.StaticInfoOf(s.iface)
	c.Assert(si.ImplicitOnCore, Equals, true)
	c.Assert(si.ImplicitOnClassic, Equals, true)
	c.Assert(si.Summary, Equals, `allows reading from hardware random number generator`)
	c.Assert(si.BaseDeclarationSlots, testutil.Contains, "hardware-random-observe")
}

func (s *HardwareRandomObserveInterfaceSuite) TestAutoConnect(c *C) {
	// FIXME: fix AutoConnect methods to use ConnectedPlug/Slot
	c.Assert(s.iface.AutoConnect(&interfaces.Plug{PlugInfo: s.plugInfo}, &interfaces.Slot{SlotInfo: s.slotInfo}), Equals, true)
}

func (s *HardwareRandomObserveInterfaceSuite) TestInterfaces(c *C) {
	c.Check(builtin.Interfaces(), testutil.DeepContains, s.iface)
}
