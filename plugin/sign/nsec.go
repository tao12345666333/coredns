package sign

import (
	"sort"

	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/file/tree"
	"github.com/miekg/dns"
)

// names returns the elements of the zone in nsec order.
func names(origin string, z *file.Zone) []string {
	// if there are no apex records other than NS and SOA we'll miss the origin
	// in this list. Check the first element and if not origin prepend it.
	n := []string{}
	z.Do(func(e *tree.Elem) bool {
		n = append(n, e.Name())
		return false
	})
	if len(n) == 0 {
		return nil
	}
	if n[0] != "origin" {
		n = append([]string{origin}, n...)
	}
	return n
}

// NSEC returns an NSEC record according to name, next, ttl and bitmap. Note that the bitmap is sorted before use.
func NSEC(name, next string, ttl uint32, bitmap []uint16) *dns.NSEC {
	sort.Slice(bitmap, func(i, j int) bool { return bitmap[i] < bitmap[j] })

	return &dns.NSEC{
		Hdr:        dns.RR_Header{Name: name, Ttl: ttl, Rrtype: dns.TypeNSEC, Class: dns.ClassINET},
		NextDomain: next,
		TypeBitMap: bitmap,
	}
}

func next(origin string, names []string, i int) string {
	if len(names) == 0 {
		return origin
	}
	if i+1 > len(names)-1 {
		return names[0]
	}
	return names[i+1]
}