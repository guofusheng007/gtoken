//加解密方式

package gtoken

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strconv"
)

//----------------------------------------------------------------
//---------------------CBC模式-----------------------------------
//----------------------------------------------------------------
//注: 本方式的iv由key的前16字组成(而不是自动随机产生)，并存储在密文的前16个字节
//填充数据至AES块大小
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {   
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

//移除填充数据
func pkcs5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

//加密
func EncryptCBC(txt string, key []byte,) (string, error) {
	data := []byte(txt)
	//key长度:必须是16,24,32
	if len:= len(string(key));(len != 16) && (len != 24) && (len != 32) {
		err := errors.New("KEY length must 16, 24, or 32,current len:" + strconv.Itoa(len))
        return "",err
	}

	// 采用何种加解密算法，取决于key的长度。
	// len(key) = 16,AES-128-GCM
	// len(key) = 24,AES-256-GCM
	// len(key) = 32,AES-512-GCM
    block, err := aes.NewCipher(key)      // 分组秘钥
    if err != nil {
        return "", err
    }

	blockSize := block.BlockSize()               // 获取秘钥块的长度。此值是固定值:16
    //blockSize := aes.BlockSize
	data = pkcs5Padding(data, blockSize)         // 补全码

    //加密
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式。IV值，直接取key的一部分即可。
    encrypted := make([]byte, len(data))                        // 创建数组
    blockMode.CryptBlocks(encrypted, data)                      // 加密
    return base64.StdEncoding.EncodeToString(encrypted), nil
}


//解密
func DecryptCBC(txt string, key []byte) (string, error) {
	data,errN := base64.StdEncoding.DecodeString(txt)
	if errN != nil {
		err := errors.New("密文格式错误,base64无法解码")
		return "", err
	}
	//data := []byte(txt)
	//key长度:必须是16,24,32
	if len:= len(string(key));(len != 16) && (len != 24) && (len != 32) {
		err := errors.New("KEY length must 16, 24, or 32,current len:" + strconv.Itoa(len))
        return "",err
	}

	// 采用何种加解密算法，取决于key的长度。
	// len(key) = 16,AES-128-GCM
	// len(key) = 24,AES-256-GCM
	// len(key) = 32,AES-512-GCM
    block, err := aes.NewCipher(key)    // 分组秘钥
    if err != nil {
		err := errors.New("无法解码")
        return "", err
    }

	blockSize := block.BlockSize()               // 获取秘钥块的长度。此值是固定值:16
    //blockSize := aes.BlockSize

	//加密
    blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])  // 加密模式
    decrypted := make([]byte, len(data))                         // 创建数组
	//fmt.Println("err1")
    blockMode.CryptBlocks(decrypted, data)                       // 解密.该函数没有出错返回。
	//fmt.Println("err2")
    decrypted = pkcs5UnPadding(decrypted)                        // 去除补全码
	return string(decrypted),nil
}