package dice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/acarl005/stripansi"
	"github.com/fy0/procs"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type deviceFile struct {
	Display      string         `json:"display"`
	Product      string         `json:"product"`
	Device       string         `json:"device"`
	Board        string         `json:"board"`
	Model        string         `json:"model"`
	FingerPrint  string         `json:"finger_print"`
	BootId       string         `json:"boot_id"`
	ProcVersion  string         `json:"proc_version"`
	Protocol     int            `json:"protocol"` // 0: Pad 1: Android 2: Watch
	IMEI         string         `json:"imei"`
	Brand        string         `json:"brand"`
	Bootloader   string         `json:"bootloader"`
	BaseBand     string         `json:"base_band"`
	SimInfo      string         `json:"sim_info"`
	OSType       string         `json:"os_type"`
	MacAddress   string         `json:"mac_address"`
	IpAddress    []int32        `json:"ip_address"`
	WifiBSSID    string         `json:"wifi_bssid"`
	WifiSSID     string         `json:"wifi_ssid"`
	ImsiMd5      string         `json:"imsi_md5"`
	AndroidId    string         `json:"android_id"`
	APN          string         `json:"apn"`
	VendorName   string         `json:"vendor_name"`
	VendorOSName string         `json:"vendor_os_name"`
	Version      *osVersionFile `json:"version"`
}

type osVersionFile struct {
	Incremental string `json:"incremental"`
	Release     string `json:"release"`
	Codename    string `json:"codename"`
	Sdk         uint32 `json:"sdk"`
}

func randomMacAddress() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "00:16:ea:ae:3c:40"
	}
	// Set the local bit
	buf[0] |= 2
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func GenerateDeviceJson(protocol int) ([]byte, error) {
	rand.Seed(time.Now().Unix())
	bootId := uuid.New()
	imei := goluhn.Generate(15) // 注意，这个imei是完全胡乱创建的，并不符合imei规则
	androidId := fmt.Sprintf("%X", rand.Uint64())

	deviceJson := deviceFile{
		Display:      "MIUI V9.5.3.0", // Rom的名字 比如 Flyme 1.1.2（魅族rom）  JWR66V（Android nexus系列原生4.3rom）
		Product:      "sagit",         // 产品名，比如这是小米6的代号
		Device:       "sagit",
		Board:        "msm8998",                                                                     // 主板:骁龙835                                                                    //
		Brand:        "Xiaomi",                                                                      // 品牌
		Model:        "MI 6",                                                                        // 型号
		Bootloader:   "unknown",                                                                     // unknown不需要改
		FingerPrint:  "Xiaomi/sagit/sagit:8.0.0/OPR1.170623.027/V9.5.3.0.OCACNFA:user/release-keys", // 指纹
		BootId:       bootId.String(),
		ProcVersion:  "Linux version 3.10.61-7254923", // 很长，后面 builder省略了
		BaseBand:     "",                              // 基带版本 4.3CPL2-... 一大堆，直接不写
		SimInfo:      "",
		OSType:       "android",
		MacAddress:   randomMacAddress(),
		IpAddress:    []int32{192, 168, rand.Int31() % 255, rand.Int31()%253 + 2}, // 192.168.x.x
		WifiBSSID:    randomMacAddress(),
		WifiSSID:     "<unknown ssid>",
		IMEI:         imei,
		AndroidId:    androidId, // 原版的 androidId和Display内容一样，我没看协议，但是按android文档上说应该是64-bit number的hex，姑且这么做
		APN:          "wifi",
		VendorName:   "MIUI", // 这个和下面一个选项(VendorOSName)都属于意义不明，找不到相似对应，不知道是啥
		VendorOSName: "xiaomi",
		Protocol:     protocol,
		Version: &osVersionFile{
			Incremental: "OCACNFA", // Build.Version.INCREMENTAL, MIUI12: V12.5.3.0.RJBCNXM
			Release:     "11",
			Codename:    "REL",
			Sdk:         29,
		},
	}

	return json.Marshal(deviceJson)
}

var defaultConfig = `
# go-cqhttp 默认配置文件

account: # 账号相关
  uin: {QQ帐号} # QQ账号
  password: {QQ密码} # 密码为空时使用扫码登录
  encrypt: false  # 是否开启密码加密
  status: 0      # 在线状态 请参考 https://docs.go-cqhttp.org/guide/config.html#在线状态
  relogin: # 重连设置
    delay: 3   # 首次重连延迟, 单位秒
    interval: 3   # 重连间隔
    max-times: 0  # 最大重连次数, 0为无限制

  # 是否使用服务器下发的新地址进行重连
  # 注意, 此设置可能导致在海外服务器上连接情况更差
  use-sso-address: true

heartbeat:
  # 心跳频率, 单位秒
  # -1 为关闭心跳
  interval: 5

message:
  # 上报数据类型
  # 可选: string,array
  post-format: string
  # 是否忽略无效的CQ码, 如果为假将原样发送
  ignore-invalid-cqcode: false
  # 是否强制分片发送消息
  # 分片发送将会带来更快的速度
  # 但是兼容性会有些问题
  force-fragment: false
  # 是否将url分片发送
  fix-url: false
  # 下载图片等请求网络代理
  proxy-rewrite: ''
  # 是否上报自身消息
  report-self-message: false
  # 移除服务端的Reply附带的At
  remove-reply-at: false
  # 为Reply附加更多信息
  extra-reply-data: false
  # 跳过 Mime 扫描, 忽略错误数据
  skip-mime-scan: false

output:
  # 日志等级 trace,debug,info,warn,error
  log-level: warn
  # 日志时效 单位天. 超过这个时间之前的日志将会被自动删除. 设置为 0 表示永久保留.
  log-aging: 15
  # 是否在每次启动时强制创建全新的文件储存日志. 为 false 的情况下将会在上次启动时创建的日志文件续写
  log-force-new: true
  # 是否启用日志颜色
  log-colorful: true
  # 是否启用 DEBUG
  debug: false # 开启调试模式

# 默认中间件锚点
default-middlewares: &default
  # 访问密钥, 强烈推荐在公网的服务器设置
  access-token: ''
  # 事件过滤器文件目录
  filter: ''
  # API限速设置
  # 该设置为全局生效
  # 原 cqhttp 虽然启用了 rate_limit 后缀, 但是基本没插件适配
  # 目前该限速设置为令牌桶算法, 请参考:
  # https://baike.baidu.com/item/%E4%BB%A4%E7%89%8C%E6%A1%B6%E7%AE%97%E6%B3%95/6597000?fr=aladdin
  rate-limit:
    enabled: false # 是否启用限速
    frequency: 1  # 令牌回复频率, 单位秒
    bucket: 1     # 令牌桶大小

database: # 数据库相关设置
  leveldb:
    # 是否启用内置leveldb数据库
    # 启用将会增加10-20MB的内存占用和一定的磁盘空间
    # 关闭将无法使用 撤回 回复 get_msg 等上下文相关功能
    enable: true

  # 媒体文件缓存， 删除此项则使用缓存文件(旧版行为)
  cache:
    image: data/image.db
    video: data/video.db

# 连接服务列表
servers:
  # 添加方式，同一连接方式可添加多个，具体配置说明请查看文档
  #- http: # http 通信
  #- ws:   # 正向 Websocket
  #- ws-reverse: # 反向 Websocket
  #- pprof: #性能分析服务器
  # 正向WS设置
  - ws:
      # 正向WS服务器监听地址
      host: 127.0.0.1
      # 正向WS服务器监听端口
      port: {WS端口}
      middlewares:
        <<: *default # 引用默认中间件
`

func GenerateConfig(qq int64, password string, port int) string {
	ret := strings.Replace(defaultConfig, "{WS端口}", fmt.Sprintf("%d", port), 1)
	ret = strings.Replace(ret, "{QQ帐号}", fmt.Sprintf("%d", qq), 1)

	password2, _ := json.Marshal(password)
	ret = strings.Replace(ret, "{QQ密码}", fmt.Sprintf("%s", string(password2)), 1)
	return ret
}

func NewGoCqhttpConnectInfoItem(account string) *ConnectInfoItem {
	conn := new(ConnectInfoItem)
	conn.Id = uuid.New().String()
	conn.Platform = "QQ"
	conn.UseInPackGoCqhttp = true
	conn.Type = "onebot"
	conn.Enable = false
	conn.RelWorkDir = "extra/go-cqhttp-qq" + account
	return conn
}

func GoCqHttpServeProcessKill(dice *Dice, conn *ConnectInfoItem) {
	defer func() {
		defer func() {
			if r := recover(); r != nil {
				dice.Logger.Error("go-cqhttp清理报错: ", r)
				// go-cqhttp 进程退出: exit status 1
			}
		}()

		if conn.UseInPackGoCqhttp {
			conn.InPackGoCqHttpLoginSuccess = false
			conn.InPackGoCqHttpQrcodeData = nil
			conn.InPackGoCqHttpRunning = false
			conn.InPackGoCqHttpQrcodeReady = false
			conn.InPackGoCqHttpNeedQrCode = false
			conn.InPackGoCqHttpLoginDeviceLockUrl = ""

			// 注意这个会panic，因此recover捕获了
			if conn.InPackGoCqHttpProcess != nil {
				conn.InPackGoCqHttpProcess.Stop()
				conn.InPackGoCqHttpProcess.Wait() // 等待进程退出，因为Stop内部是Kill，这是不等待的
			}
		}
	}()
}

func GoCqHttpServeRemoveSessionToken(dice *Dice, conn *ConnectInfoItem) {
	workDir := filepath.Join(dice.BaseConfig.DataDir, conn.RelWorkDir)
	if _, err := os.Stat(filepath.Join(workDir, "session.token")); err == nil {
		os.Remove(filepath.Join(workDir, "session.token"))
	}
}

func GoCqHttpServe(dice *Dice, conn *ConnectInfoItem, password string, protocol int, isAsyncRun bool) {
	workDir := filepath.Join(dice.BaseConfig.DataDir, conn.RelWorkDir)
	os.MkdirAll(workDir, 0644)

	qrcodeFile := filepath.Join(workDir, "qrcode.png")
	deviceFilePath := filepath.Join(workDir, "device.json")
	configFilePath := filepath.Join(workDir, "config.yml")
	if _, err := os.Stat(qrcodeFile); err == nil {
		// 如果已经存在二维码文件，将其删除
		os.Remove(qrcodeFile)
		fmt.Println("删除已存在的二维码文件")
	}

	//if _, err := os.Stat(filepath.Join(workDir, "session.token")); errors.Is(err, os.ErrNotExist) {
	if !conn.InPackGoCqHttpLoginSucceeded {
		// 并未登录成功，删除记录文件
		os.Remove(configFilePath)
		os.Remove(deviceFilePath)
	}

	// 创建设备配置文件
	if _, err := os.Stat(deviceFilePath); errors.Is(err, os.ErrNotExist) {
		deviceInfo, err := GenerateDeviceJson(protocol)
		if err == nil {
			ioutil.WriteFile(deviceFilePath, deviceInfo, 0644)
		}
	} else {
		fmt.Println("设备文件已存在，跳过")
	}

	// 创建配置文件
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		// 如果不存在 config.yml 那么启动一次，让它自动生成
		// 改为：如果不存在，帮他创建
		p, _ := GetRandomFreePort()
		conn.ConnectUrl = fmt.Sprintf("ws://localhost:%d", p)
		c := GenerateConfig(conn.UserId, password, p)
		ioutil.WriteFile(configFilePath, []byte(c), 0644)
	}

	// 启动客户端
	wd, _ := os.Getwd()
	gocqhttpExePath, _ := filepath.Abs(filepath.Join(wd, "go-cqhttp/go-cqhttp"))
	gocqhttpExePath = strings.Replace(gocqhttpExePath, "\\", "/", -1) // windows平台需要这个替换

	p := procs.NewProcess(gocqhttpExePath + " faststart")
	p.Dir = workDir

	chQrCode := make(chan int, 1)
	p.OutputHandler = func(line string) string {
		// 请使用手机QQ扫描二维码 (qrcode.png) :
		if strings.Contains(line, "qrcode.png") {
			chQrCode <- 1
		}
		if strings.Contains(line, "CQ WebSocket 服务器已启动") {
			// CQ WebSocket 服务器已启动
			// 登录成功 欢迎使用
			conn.InPackGoCqHttpLoginSuccess = true
			conn.InPackGoCqHttpLoginSucceeded = true
			conn.Enable = true
			conn.State = 1
			dice.Logger.Infof("gocqhttp登录成功，帐号: <%s>(%d)", conn.Nickname, conn.UserId)

			go DiceServe(dice, conn)
		}

		if strings.Contains(line, "账号已开启设备锁，请前往") {
			re := regexp.MustCompile(`-> (.+?) <-`)
			m := re.FindStringSubmatch(line)
			if len(m) > 0 {
				// 设备锁流程，因为需要重新登录，进行一个“已成功登录过”的标记，这样配置文件不会被删除
				conn.InPackGoCqHttpLoginSucceeded = true
				conn.InPackGoCqHttpLoginDeviceLockUrl = m[1]
			}
		}

		if strings.Contains(line, "请使用手机QQ扫描二维码以继续登录") {
			conn.InPackGoCqHttpNeedQrCode = true
		}

		if conn.InPackGoCqHttpLoginSuccess == false || strings.Contains(line, "风控") || strings.Contains(line, "WARNING") || strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
			//  [WARNING]: 登录需要滑条验证码, 请使用手机QQ扫描二维码以继续登录
			if conn.InPackGoCqHttpLoginSuccess {
				dice.Logger.Infof("onebot | %s", stripansi.Strip(line))
			} else {
				fmt.Printf("onebot | %s\n", line)

				// error 之类错误无条件警告
				if strings.Contains(line, "WARNING") || strings.Contains(line, "ERROR") || strings.Contains(line, "FATAL") {
					dice.Logger.Infof("onebot | %s", stripansi.Strip(line))
				}
			}
		}
		return line
	}

	go func() {
		<-chQrCode
		if _, err := os.Stat(qrcodeFile); err == nil {
			fmt.Println("二维码已经就绪")
			fmt.Println("如控制台二维码不好扫描，可以手动打开go-cqhttp目录下qrcode.png")
			qrdata, err := ioutil.ReadFile(qrcodeFile)
			if err == nil {
				conn.InPackGoCqHttpQrcodeData = qrdata
			}
			conn.InPackGoCqHttpQrcodeReady = true
		}
	}()

	run := func() {
		conn.InPackGoCqHttpRunning = true
		conn.InPackGoCqHttpProcess = p
		err := p.Run()
		GoCqHttpServeProcessKill(dice, conn)
		conn.InPackGoCqHttpRunning = false
		if err != nil {
			dice.Logger.Info("go-cqhttp 进程退出: ", err)
		} else {
			dice.Logger.Info("go-cqhttp 进程退出")
		}
	}

	if isAsyncRun {
		go run()
	} else {
		run()
	}
}

func DiceServe(d *Dice, conn *ConnectInfoItem) {
	if !conn.DiceServing {
		conn.DiceServing = true
	} else {
		return
	}

	var _index int
	for index, a := range d.ImSession.Conns {
		if a == conn {
			_index = index
		}
	}

	checkQuit := func() bool {
		if !conn.DiceServing {
			// 退出连接
			d.Logger.Infof("检测到连接删除，不再进行此onebot服务的重连: <%s>(%d)", conn.Nickname, conn.UserId)
			return true
		}
		return false
	}

	for {
		if checkQuit() {
			break
		}

		// 骰子开始连接
		d.Logger.Infof("开始连接 onebot 服务，帐号 <%s>(%d)", conn.Nickname, conn.UserId)
		ret := d.ImSession.Serve(_index)

		if ret == 0 {
			break
		}

		if checkQuit() {
			break
		}

		d.Logger.Infof("onebot 连接中断，将在15秒后重新连接，帐号 <%s>(%d)", conn.Nickname, conn.UserId)
		time.Sleep(time.Duration(15 * time.Second))
	}
}