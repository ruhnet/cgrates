/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/
package engine

import (
	"fmt"
	"testing"

	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/utils"
	"github.com/cgrates/rpcclient"
	"github.com/nyaruka/phonenumbers"
)

func TestDynamicDpFieldAsInterface(t *testing.T) {
	cfg := config.NewDefaultCGRConfig()
	ms := utils.MapStorage{}
	dDp := newDynamicDP([]string{}, []string{utils.ConcatenatedKey(utils.MetaInternal, utils.StatSConnsCfg)}, []string{}, "cgrates.org", ms)
	clientconn := make(chan rpcclient.ClientConnector, 1)
	clientconn <- &ccMock{
		calls: map[string]func(args interface{}, reply interface{}) error{
			utils.StatSv1GetQueueFloatMetrics: func(args, reply interface{}) error {
				rpl := &map[string]float64{
					"stat1": 31,
				}
				*reply.(*map[string]float64) = *rpl
				return nil
			},
		},
	}
	connMgr := NewConnManager(cfg, map[string]chan rpcclient.ClientConnector{
		utils.ConcatenatedKey(utils.MetaInternal, utils.StatSConnsCfg): clientconn,
	})
	SetConnManager(connMgr)
	if _, err := dDp.fieldAsInterface([]string{utils.MetaStats, "val", "val3"}); err == nil || err != utils.ErrNotFound {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaLibPhoneNumber, "+402552663", "val3"}); err != nil {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaLibPhoneNumber, "+402552663", "val3"}); err != nil {
		t.Error(err)
	} else if _, err := dDp.fieldAsInterface([]string{utils.MetaAsm, "+402552663", "val3"}); err == nil || err != utils.ErrNotFound {
		t.Error(err)
	}

}

func TestDpLibPhoneNumber(t *testing.T) {
	num, err := phonenumbers.ParseAndKeepRawInput("+3554735474", utils.EmptyString)
	if err != nil {
		t.Error(err)
	}
	dDP := &libphonenumberDP{
		pNumber: num,
		cache:   utils.MapStorage{},
	}
	dDP.setDefaultFields()
	if _, err := dDP.fieldAsInterface([]string{}); err == nil || err.Error() != fmt.Sprintf("invalid field path <%+v> for libphonenumberDP", []string{}) {
		t.Error(err)
	}
	val, err := dDP.fieldAsInterface([]string{"CountryCode"})
	if err != nil {
		t.Error(err)
	}
	if err != nil {
		t.Error(err)
	}
	exp := int32(355)
	if nationalNumber, cancast := val.(int32); !cancast {
		t.Error("can't convert")
	} else if nationalNumber != exp {
		t.Errorf("expected %v,received %v", exp, nationalNumber)
	}

}