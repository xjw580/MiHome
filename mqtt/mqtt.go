package mqtt

import "C"
import (
	"bemfa-demo/systemUtil"
	"bemfa-demo/wyyUtil"
	"bufio"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-toast/toast"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	defaultConfig  = [...]string{"cloudmusicPath=", "key=", "topicName=", "iconPath="}
	key            = "key"
	cloudmusicPath = "cloudmusicPath"
	topic          = "topicName"
	pubSet         = "/set"
	pubUp          = "/up"
	iconPath       = "iconPath"
)

const (
	topicType   = "005"
	server      = "bemfa.com"
	port        = "9501"
	configPath  = "config.properties"
	programName = "MiHome"
)

var (
	client      MQTT.Client
	responseMap = map[string]func(client MQTT.Client, message MQTT.Message){
		"on": func(client MQTT.Client, message MQTT.Message) {
			pids, _ := systemUtil.GetProcessIDsByName(wyyUtil.CloudmusicName)
			if pids == nil || len(pids) == 0 {
				go func() {
					cmd := exec.Command(cloudmusicPath)
					_, err := cmd.CombinedOutput()
					if err != nil {
						log.Fatalln(err)
					}
				}()
				time.Sleep(6 * time.Second)
			}

			if !wyyUtil.GetPlayStatus() {
				wyyUtil.ChangePlayStatus()
			}
		},
		"off": func(client MQTT.Client, message MQTT.Message) {
			if wyyUtil.GetPlayStatus() {
				wyyUtil.ChangePlayStatus()
			}
		},
	}
)

func init() {
	if !loadConfig() {
		notice(programName, "error", "读取配置文件异常", iconPath, "", "")
		log.Println("读取配置文件异常")
		os.Exit(-1)
	}
	if key == "" {
		notice(programName, "error", "Failed to connect to MQTT broker", iconPath, "", "")
		return
	}
	// 定义MQTT连接选项
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://" + server + ":" + port) // 替换为你的MQTT服务器地址
	opts.SetClientID(key)
	opts.SetUsername("") // 如果需要用户名密码认证，设置用户名
	opts.SetPassword("") // 设置密码
	opts.DefaultPublishHandler = messagePubHandler
	opts.OnConnect = connectSuccessHandler
	opts.OnConnectionLost = connectLostHandler

	// 创建MQTT客户端
	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		notice(programName, "error", "Failed to connect to MQTT broker", iconPath, "", "")
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

}

func notice(appID, title, msg, icoPath, btnText, btnURL string) {
	var notification toast.Notification
	if btnText == "" {
		notification = toast.Notification{
			AppID:   appID,
			Title:   title,
			Message: msg,
			Icon:    icoPath,
		}
	} else {
		notification = toast.Notification{
			AppID:   appID,
			Title:   title,
			Message: msg,
			Icon:    icoPath,
			Actions: []toast.Action{
				{"protocol", btnText, btnURL},
			},
		}
	}

	err := notification.Push()
	if err != nil {
		panic(err)
	}
}

func loadConfig() bool {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 文件不存在，创建并写入数据
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return false
		}
		defer file.Close()

		// 创建一个带缓冲的写入器
		writer := bufio.NewWriter(file)

		for _, s := range defaultConfig {
			// 将数据写入文件
			_, err = writer.WriteString(s + "\n")
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return false
			}
		}

		// 刷新缓冲区，确保数据被写入文件
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer:", err)
			return false
		}

		fmt.Println("Data written to file successfully.")
		return false
	} else {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return false
		}
		// 创建一个 map 用于存储属性
		properties := make(map[string]string)

		// 创建一个 Scanner 对象来逐行读取文件内容
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			// 忽略空行和注释行
			if len(line) > 0 && !strings.HasPrefix(line, "#") {
				// 分割属性名和属性值
				parts := strings.Split(line, "=")
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					properties[key] = value
				}
			}
		}

		// 检查扫描过程中是否出现错误
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			return false
		}

		key = properties[key]
		cloudmusicPath = properties[cloudmusicPath]
		topic = properties[topic] + topicType
		pubSet = topic + topicType
		pubUp = topic + topicType
		iconPath = properties[iconPath]
		// 打印 properties map
		for key, value := range properties {
			fmt.Printf("%s = %s\n", key, value)
		}
		return true
	}
}

func Launch() {
	if token := client.Subscribe(topic, 1, func(client MQTT.Client, message MQTT.Message) {
		msg := string(message.Payload())
		fmt.Println("received: ", msg)
		f := responseMap[msg]
		if f == nil {
			split := strings.Split(msg, "#")
			l := len(split)
			if l == 6 {
				wyyUtil.ChangeVoice(split[5] == "1")
			} else if l == 5 {
				wyyUtil.ChangeMusic(split[4] == "1")
			}

		} else {
			f(client, message)
		}
		message.Ack()
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	} else {
		sprintf := fmt.Sprintf("Subscribed to topic: %s\n", topic)
		notice(programName, "successful", sprintf, iconPath, "", "")
	}
}

// 发布消息
func PublicMsg(text string) {
	if token := client.Publish(pubSet, 1, false, text); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to publish message: %v", token.Error())
	}
	fmt.Printf("Published set message: %s\n", text)
}

// 发布消息的处理器
var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// 连接成功的处理器
var connectSuccessHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	fmt.Println("Connected to MQTT broker")
}

// 连接丢失的处理器
var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("reboot...")
	notice(programName, "successful", "reboot...", iconPath, "", "")
	go func() {
		exec.Command(execPath).Start()
	}()
	time.Sleep(250 * time.Millisecond)
	os.Exit(0)
}
