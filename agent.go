package devcontainershell

import (
	"context"
	"encoding/json"
	"net/netip"
	"os"
	"time"

	"github.com/yskszk63/netlink-list-listens"
)

type ListenEvent struct {
	Type string `json:"type"`
	IPv6 bool   `json:"ipv6"`
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
}

func update(m *map[netip.AddrPort]struct{}) ([]netip.AddrPort, []netip.AddrPort, error) {
	l, err := netlinklistlistens.ListListens()
	if err != nil {
		return nil, nil, err
	}

	old := *m
	new := make(map[netip.AddrPort]struct{})
	*m = new
	add := make([]netip.AddrPort, 0)
	rm := make([]netip.AddrPort, 0)

	for _, addr := range l {
		new[addr] = struct{}{}

		_, exists := old[addr]
		if !exists {
			add = append(add, addr)
			continue
		}
		delete(old, addr)
	}

	for k := range old {
		rm = append(rm, k)
	}

	return add, rm, nil
}

func WatchListens(cx context.Context) error {
	m := make(map[netip.AddrPort]struct{})

	d := time.Second * 1
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	enc := json.NewEncoder(os.Stdout)

	for {
		add, rm, err := update(&m)
		if err != nil {
			return err
		}

		for _, a := range add {
			msg := ListenEvent{
				Type: "ADD",
				IPv6: a.Addr().Is6(),
				IP:   a.Addr().String(),
				Port: a.Port(),
			}
			if err := enc.Encode(msg); err != nil {
				return err
			}
		}

		for _, r := range rm {
			msg := ListenEvent{
				Type: "REMOVE",
				IPv6: r.Addr().Is6(),
				IP:   r.Addr().String(),
				Port: r.Port(),
			}
			if err := enc.Encode(msg); err != nil {
				return err
			}
		}

		select {
		case <-cx.Done():
			return nil
		case <-ticker.C:
		}
	}
}
