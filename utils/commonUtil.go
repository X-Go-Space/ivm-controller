package utils

import "github.com/duke-git/lancet/v2/random"

const (
    QRcodeInit = "init"
    QRcodeScan = "scan"
    QRcodeConfirm = "confirm"
    QRcodeExpired = "expired"
)

func RandID()string {
    return random.RandNumeralOrLetter(16)
}

func GenerateResourceId(ID string) string{
    return "RESOURCE##"+ID
}
