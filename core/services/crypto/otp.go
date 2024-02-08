package crypto

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func ValidateOTPToken(otp string) bool {
	secret := os.Getenv("CORE_OTP_SECRET")
	return totp.Validate(otp, secret)
}

func GetCurrentOTP() string {
	secret := os.Getenv("CORE_OTP_SECRET")
	code, _ := totp.GenerateCode(secret, time.Now())
	return code
}

func GenesisQRCode() {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "Itsramp",
		AccountName: "bruno@bruno.adm.br",
		SecretSize:  64,
	})
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, _ := key.Image(200, 200)
	png.Encode(&buf, img)
	// display the QR code to the user.
	display(key, buf.Bytes())

}

func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	ioutil.WriteFile("qr-code.png", data, 0644)
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}
