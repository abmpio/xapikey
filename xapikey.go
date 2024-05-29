package xapikey

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/abmpio/entity"
	uuid "github.com/satori/go.uuid"
)

// AKSK 结构体定义了 Access Key（AK）和 Secret Key（SK）的结构
type Aksk struct {
	entity.EntityWithUser `bson:",inline"`
	// 所属app
	App string `json:"app" bson:"app"`

	// Access Key（AK）
	AccessKey string `json:"accessKey" bson:"accessKey"`
	// Secret Key（SK）
	SecretKey      string    `json:"secretKey" bson:"secretKey"`
	CreatedTime    time.Time `json:"createdTime" bson:"createdTime"`
	ExpirationTime time.Time `json:"expirationTime" bson:"expirationTime"`
	// 状态：启用、禁用等
	Status bool `json:"status" bson:"status"`
	// ip白名单,多个ip以;隔开
	IpWhitelist string `json:"ipWhitelist" bson:"ipWhitelist"`

	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

type IAkskService interface {
	entity.IEntityService[Aksk]

	//根据app与ak来查找列表
	FindByAk(app string, ak string) ([]*Aksk, error)
}

// 生成 AK/SK 的函数
func GenerateAKSK() (string, string) {
	// 生成随机的 AK（使用 UUID）
	ak := uuid.NewV4().String()

	// 生成随机的 SK（使用当前时间作为盐值）
	salt := fmt.Sprintf("%d", time.Now().UnixNano())
	sk := hashSaltedSK(ak, salt)

	return ak, sk
}

// 对 SK 进行加盐和哈希处理
func hashSaltedSK(sk, salt string) string {
	// 使用 HMAC-SHA256 算法对 SK 进行加盐和哈希处理
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(sk))
	hashedSK := h.Sum(nil)

	// 返回经过 base64 编码的哈希后的 SK 和盐值
	return base64.StdEncoding.EncodeToString(hashedSK)
}
