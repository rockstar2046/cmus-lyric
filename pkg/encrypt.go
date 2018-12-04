// CODE FROM GITHUB.COM
package pkg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"math/big"
	"math/rand"
)

const modulus = "00e0b509f6259df8642dbc35662901477df22677ec152b5ff68ace615bb7b725152b3ab17a876aea8a5aa76d2e417629ec4ee341f56135fccf695280104e0312ecbda92557c93870114af6c9d05c4f7f0c3685b7a46bee255932575cce10b424d813cfe4875d3e82047b97ddef52741d546b8e289dc6935b3ece0462db0a22b8e7"
const nonce = "0CoJUm6Qyw8W8jud"
const pubKey = "010001"
const keys = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/"
const iv = "0102030405060708"

// EncParams 传入参数 得到加密后的参数和一个也被加密的秘钥
func EncParams(param string) (string, string, error) {
	// 创建 key
	secKey := createSecretKey(16)

	aes1, err1 := aesEncrypt(param, nonce)

	// 第一次加密 使用固定的 nonce
	if err1 != nil {
		return "", "", err1
	}
	aes2, err2 := aesEncrypt(aes1, secKey)
	// 第二次加密 使用创建的 key
	if err2 != nil {
		return "", "", err2
	}
	// 得到 加密好的 param 以及 加密好的key
	return aes2, rsaEncrypt(secKey, pubKey, modulus), nil
}

// 创建指定长度的key
func createSecretKey(size int) string {
	// 也就是从 a~9 以及 +/ 中随机拿出指定数量的字符拼成一个 size长的key
	rs := ""
	for i := 0; i < size; i++ {
		pos := rand.Intn(len([]rune(keys)))
		rs += keys[pos : pos+1]
	}
	return rs
}

// 通过 CBC模式的AES加密 用 sKey 加密 sSrc
func aesEncrypt(sSrc string, sKey string) (string, error) {
	iv := []byte(iv)
	block, err := aes.NewCipher([]byte(sKey))
	if err != nil {
		return "", err
	}
	//需要padding的数目
	padding := block.BlockSize() - len([]byte(sSrc))%block.BlockSize()
	//只要少于256就能放到一个byte中，默认的blockSize=16(即采用16*8=128, AES-128长的密钥)
	//最少填充1个byte，如果原文刚好是blocksize的整数倍，则再填充一个blocksize

	src := append([]byte(sSrc), bytes.Repeat([]byte{byte(padding)}, padding)...)

	model := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(src))
	model.CryptBlocks(cipherText, src)
	// 最后使用base64编码输出
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// 将 key 也加密
func rsaEncrypt(key string, pubKey string, modulus string) string {
	// 倒序 key
	rKey := ""
	for i := len(key) - 1; i >= 0; i-- {
		rKey += key[i : i+1]
	}
	// 将 key 转 ascii 编码 然后转成 16 进制字符串
	hexRKey := ""
	for _, char := range []rune(rKey) {
		hexRKey += fmt.Sprintf("%x", int(char))
	}
	// 将 16进制 的 三个参数 转为10进制的 bigint
	bigRKey, _ := big.NewInt(0).SetString(hexRKey, 16)
	bigPubKey, _ := big.NewInt(0).SetString(pubKey, 16)
	bigModulus, _ := big.NewInt(0).SetString(modulus, 16)
	// 执行幂乘取模运算得到最终的bigint结果
	bigRs := bigRKey.Exp(bigRKey, bigPubKey, bigModulus)
	// 将结果转为 16进制字符串
	hexRs := fmt.Sprintf("%x", bigRs)
	// 可能存在不满256位的情况，要在前面补0补满256位
	return addPadding(hexRs, modulus)
}

// 补0步骤
func addPadding(encText string, modulus string) string {
	ml := len(modulus)
	for i := 0; ml > 0 && modulus[i:i+1] == "0"; i++ {
		ml--
	}
	num := ml - len(encText)
	prefix := ""
	for i := 0; i < num; i++ {
		prefix += "0"
	}
	return prefix + encText
}
