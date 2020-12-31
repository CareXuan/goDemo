package api

import "github.com/skip2/go-qrcode"

func getQrCode(qrUrl string, picPath string) {
	qrcode.WriteFile(qrUrl, qrcode.Medium, 256, picPath)
}
