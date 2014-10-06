package lldp

func (this LinkDiscoveryTLV) Size() int {
	tl := this.TypeAndLen()
	return int(tl&0x01FF + 2)
}

func (this LinkDiscoveryTLV) SetSize(s int) {
	this.MergeTypeAndLen(this.Type(), uint16(s-2))
}

func (this LinkDiscoveryTLV) MergeTypeAndLen(t uint8, l uint16) {
	this.SetTypeAndLen(uint16(uint16(t&0x07F)<<9 | l&0x1FF))
}

func (this LinkDiscoveryTLV) Type() uint8 {
	return uint8((this.TypeAndLen() & 0xFF00) >> 9)
}

func (this LinkDiscoveryTLV) SetType(t uint8) {
	this.MergeTypeAndLen(t, uint16(this.Size()-2))
}
