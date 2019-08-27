package main

import (
	"testing"

	"golang.org/x/net/publicsuffix"
)

func TestDomainName(t *testing.T) {
	s1, _ := publicsuffix.EffectiveTLDPlusOne("http://wwww.baidu.com")
	t.Log(s1)
	s2, _ := publicsuffix.EffectiveTLDPlusOne("http://pan.baidu.com")
	t.Log(s2)
	eTLD, icann := publicsuffix.PublicSuffix("http://wwww.baidu.com")
	t.Log(eTLD, icann)
}
