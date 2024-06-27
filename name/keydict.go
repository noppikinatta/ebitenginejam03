package name

var txtimgKeyMapping map[string]string = map[string]string{
	TextKeyEquip1Laser:      ImgKeyEquip1Laser,
	TextKeyEquip2Missile:    ImgKeyEquip2Missile,
	TextKeyEquip3Harakiri:   ImgKeyEquip3Harakiri,
	TextKeyEquip4Barrier:    ImgKeyEquip4Barrier,
	TextKeyEquip5Armor:      ImgKeyEquip5Armor,
	TextKeyEquip6Exhaust:    ImgKeyEquip6Exhaust,
	TextKeyEquip7Stonehenge: ImgKeyEquip7Stonehenge,
	TextKeyEquip8Sushibar:   ImgKeyEquip8Sushibar,
	TextKeyEquip9Operahouse: ImgKeyEquip9Operahouse,
	TextKeyManager1:         ImgKeyManager1,
	TextKeyManager2:         ImgKeyManager2,
	TextKeyManager3:         ImgKeyManager3,
	TextKeyVendor1:          ImgKeyVendor1,
	TextKeyVendor2:          ImgKeyVendor2,
	TextKeyVendor3:          ImgKeyVendor3,
}

func ImgKey(txtKey string) string {
	k, ok := txtimgKeyMapping[txtKey]
	if ok {
		return k
	}

	// fallback image may be drawn
	return "fallback"
}
