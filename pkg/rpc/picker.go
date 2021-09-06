/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rpc

import (
	"context"
	"log"

	"github.com/serialx/hashring"
	"google.golang.org/grpc/balancer"
)

const (
	Key2taskId = "task_id"
)

var (
	_ balancer.Picker = &d7yPicker{}
)

type pickResult struct {
	Ctx context.Context
	SC  balancer.SubConn
}

func NewD7yPicker(subConns map[string]balancer.SubConn) *d7yPicker {
	addrs := make([]string, 0)
	for addr := range subConns {
		addrs = append(addrs, addr)
	}
	return &d7yPicker{
		subConns:   subConns,
		hashRing:   hashring.New(addrs),
		needReport: false,
	}
}

func NewD7yReportingPicker(subConns map[string]balancer.SubConn, reportChan chan<- pickResult) *d7yPicker {
	addrs := make([]string, 0)
	for addr := range subConns {
		addrs = append(addrs, addr)
	}
	return &d7yPicker{
		subConns:   subConns,
		hashRing:   hashring.New(addrs),
		needReport: true,
		reportChan: reportChan,
	}
}

type d7yPicker struct {
	subConns   map[string]balancer.SubConn // address string -> balancer.SubConn
	hashRing   *hashring.HashRing
	needReport bool
	reportChan chan<- pickResult
}

func (p *d7yPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var ret balancer.PickResult
	key, ok := info.Ctx.Value(Key2taskId).(string)
	if !ok {
		// for keepAlive
		key = info.FullMethodName
	}
	log.Printf("pick for %s\n", key)
	if targetAddr, ok := p.hashRing.GetNode(key); ok {
		ret.SubConn = p.subConns[targetAddr]
		if p.needReport {
			p.reportChan <- pickResult{Ctx: info.Ctx, SC: ret.SubConn}
		}
	}
	if ret.SubConn == nil {
		return ret, balancer.ErrNoSubConnAvailable
	}
	return ret, nil
}
