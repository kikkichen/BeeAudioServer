package utils

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const FUNNY_SOUL = "2330418701,3194101060,2102631590,2448497652,5262258447,5281293069,5204999492,5499154875,6610518564,6621497083,5354035546,5657825155,5531311536,1758307583,7609058227,6027232453,3915151208,2844929434,5563152336,6494027596,7785070452,6733296921,1004745095,1089477917,2015145765,1642484781,1761934994,3282943757,1242846212,6859251176,6835553691,5329047588,3769758905,2822676681,3757408904,2473790195,5066348687,5096341817,7487904465,5512767483,5312999796,5486349335,2825485270,2671441390,5647985811,5210717813,5992522323,5717974866,3268990791,7570137873,6905945547,1559505935,6305386541,2502184780,6296664021,6598737815,3288207015,6503525169,2994814360,7734957731,1791218150,5145707462,3178063880,3470168084,2082348875,7388541685,3032064021,5743889077,6116796498"

/* 字符串ID 转 uint64 类型 */
func StringParseToUint64(number string) uint64 {
	uintId, err := strconv.ParseUint(number, 10, 64)
	if err != nil {
		return 0
	}
	return uintId
}

var DefaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-*#")
var DefaultUpperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var DefaultNumber = []rune("1234567890")

/*
 *	生成指定长度的随机字符串
 *	@params	n	指定字符串长度
 *	@parmas	allowedChars	可选参数，字符范围集合
 */
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = DefaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}

	return string(b)
}

/*
 *	生成指定长度的随机字符串
 *	@params	n	指定字符串长度
 *	@parmas	allowedChars	可选参数，字符范围集合
 */
func RandomNumberString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = DefaultNumber
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[n.Int64()]
	}

	return string(b)
}

/*
 *	生成md5加密后的字符串
 *	@params	text 待转为加密字符串的字符串
 */
func GenerateStringByMD5(text string) string {
	src_code := md5.Sum([]byte(text))
	code := fmt.Sprintf("%x", src_code)
	return string(code)
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
