package update

import (
	"context"

	"github.com/Autumn-27/ScopeSentry/internal/constants"
	"github.com/Autumn-27/ScopeSentry/internal/database/mongodb"
	"github.com/Autumn-27/ScopeSentry/internal/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// builtinPluginTexts 는 기본 제공 플러그인의 설명문이 메시지 키로 바뀌기 전에
// DB 에 저장되어 있던 중국어 문장이다. 키는 플러그인 hash 다.
//
// 이 표가 있어야 하는 이유: 기존 설치본의 plugins 컬렉션에는 중국어 문장이 들어
// 있고, 새 코드는 같은 자리에서 메시지 키를 기대한다. 저장된 값이 아래 문장과
// 정확히 같을 때만 키로 바꿔서, 사용자가 손댄 값은 건드리지 않는다.
var builtinPluginTexts = map[string]struct {
	Help         string
	Introduction string
}{
	"80718cc3fcb4827d942e6300184707e2": {"无需参数", "web指纹识别"},
	"3a0d994a12305cb15a5cb7104d819623": {"-cdncheck 是否开启cdn检测 -screenshot 是否开启截图，默认关闭,开启需要安装chromium -tlsprobe 从tls信息发送http探测默认true", "资产测绘"},
	"920546788addc6d29ea63e4a314a1b85": {"-d 目录扫描字典 -t 扫描并发限制", "目录扫描"},
	"648a6f49eed57b1737ac702e02985b00": {"无需参数", "端口指纹识别"},
	"66b4ddeb983387df2b7ee7726653874d": {"-port 端口扫描范围 -b 端口扫描并发数量  -t 超时时间", "端口存活扫描"},
	"9b91e0f18ac9043ec9fe250a39b4a2d9": {"无需参数", "检测是否cdn，跳过cdn的端口扫描"},
	"d60ba73c70aac430a0a54e796e7e19b8": {"-t 扫描线程  -timeout 超时时间 -max-time 最大等待时间", "子域名扫描"},
	"e8f55f5e0e9f4af1ca40eb19048b8c82": {"-subfile 子域名字典 -et 最长运行时间(分钟)", "子域名爆破"},
	"c0c71c101271f38b8be1767f3626d291": {"无需参数", "子域名接管检测"},
	"ef244b3462744dad3040f9dcf3194eb1": {"无需参数", "url扫描从Waybackarchive、Alienvault、Commoncrawl获取历史url"},
	"9669d0dcc52a5ca6dbbe580ffc99c364": {"-t 并发数 -timeout 超时时间 -et 最长运行时间(分钟) -rs 读取页面大小(MB)", "url爬取"},
	"2949994c04a4e124b9c98383489510f0": {"无需参数", "敏感信息泄露检测"},
	"e52b8b16d49912ca564c22319c495403": {"无需参数", "页面监控，将所有url放入页面监控的计划任务中"},
	"1aa212b9578dc3fb1409ee8de8ed005e": {"-pdf 开启pdf检测 -exclude 排除提取的规则(name1,name2) -verify 是否进行验证（验证通过再统计结果）", "trufflehog密钥提取，如果设置了排除规则，需要重新安装才能重新启用被排除的规则"},
	"ed93b8af6b72fe54a60efdb932cf6fbc": {"参考官方支持t, s, es, tags, etags, rl, rld, bs, c, hbs, headc, jsc, pc, prc参数", "漏洞扫描"},
	"4b292861d3228af0e4da8e7ef979497c": {"参考插件市场说明", "爬虫"},
}

// Update21 은 기본 제공 플러그인의 help / introduction 을 메시지 키로 바꾼다.
// 멱등하므로 매번 호출해도 된다. 이미 키가 들어 있는 행은 조건에 걸리지 않는다.
func Update21() {
	coll := mongodb.DB.Collection("plugins")

	wanted := make(map[string]struct{ Help, Introduction string }, len(constants.Plugins))
	for _, p := range constants.Plugins {
		wanted[p.Hash] = struct{ Help, Introduction string }{p.Help, p.Introduction}
	}

	migrated := 0
	for hash, old := range builtinPluginTexts {
		keys, ok := wanted[hash]
		if !ok {
			continue
		}

		res, err := coll.UpdateOne(
			context.Background(),
			bson.M{"hash": hash, "help": old.Help, "introduction": old.Introduction},
			bson.M{"$set": bson.M{"help": keys.Help, "introduction": keys.Introduction}},
		)
		if err != nil {
			logger.Error("failed to migrate built-in plugin texts", zap.String("hash", hash), zap.Error(err))
			continue
		}
		migrated += int(res.ModifiedCount)
	}

	if migrated > 0 {
		logger.Info("migrated built-in plugin descriptions to message keys", zap.Int("count", migrated))
	}
}
