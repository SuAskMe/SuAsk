package utility

import (
	"golang.org/x/crypto/bcrypt"
)

// -----------------------------------------------------------------------------
// 密码哈希策略
//
// 新账号 / 修改密码 → bcrypt（cost=10）。
// 历史账号仍然是 md5(md5(pw)+md5(salt))，登录时透明升级到 bcrypt，
// 第二次登录往后 DB 里就只剩 bcrypt 了。
// -----------------------------------------------------------------------------

// bcryptCost 选 10 是在安全性和延迟之间的常见折中。
// 一次校验在现代服务器 CPU 上大约 60-80ms，登录场景完全够用。
const bcryptCost = 10

// HashPassword 生成 bcrypt 密码哈希。bcrypt 自带盐，所以不再需要外部 salt 参数。
func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

// IsLegacyHash 判断 DB 里存的是否还是老的 MD5 格式。
// 老格式长度固定 32，新格式以 $2a$ / $2b$ / $2y$ 开头且长度 60。
// 把判定集中在这里，避免各个调用点重复写魔法字符串。
func IsLegacyHash(hash string) bool {
	if len(hash) != 32 {
		return false
	}
	for _, c := range hash {
		// 任何非十六进制字符都不可能是老的 md5 输出
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// VerifyPassword 对明文密码做校验。
//
//   - 新 bcrypt 哈希：走 bcrypt.CompareHashAndPassword。
//   - 旧 MD5 哈希：复用 EncryptPassword 的算法，用常量时间比较。
//     调用方拿到 shouldUpgrade=true 时应该重新算 bcrypt 并写回 DB。
//
// 返回值：
//
//	match       密码是否正确
//	needUpgrade 密码正确，但底层还在用旧算法，建议升级
//	err         只在 bcrypt 库内部错误时非空；密码错不算 err
func VerifyPassword(stored, salt, password string) (match, needUpgrade bool, err error) {
	if IsLegacyHash(stored) {
		if constantTimeEqualStr(EncryptPassword(password, salt), stored) {
			return true, true, nil
		}
		return false, false, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored), []byte(password))
	if err == nil {
		return true, false, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, false, nil
	}
	// 其它错误（例如 hash 格式异常），返回给上层记录日志
	return false, false, err
}

// constantTimeEqualStr 对字符串做恒定时间比较，避免 md5 场景下通过时间差做猜测。
// 标准库有 subtle.ConstantTimeCompare，但那个只认 []byte；自己包一层保持 API 清爽。
func constantTimeEqualStr(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := 0; i < len(a); i++ {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}
