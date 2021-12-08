#include <ESP8266WiFi.h>
#define led 2
const char *ssid = "Nuclear Fusion Reactor";//网络的ssid
const char *password = "********";//wifi密码
const char *host = "192.168.43.166";//Server服务端的IP
WiFiClient client;
const int tcpPort = 8266;//Server服务端的端口号
static String comdata = "";
static String val = "";

void setup() {
    Serial.begin(9600);
    pinMode(led, OUTPUT);
    delay(10);
    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);
    //检测是否成功连接到目标网络，未连接则阻塞。
    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
    }
    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
    Serial.println("*");
}

void loop() {
    listenSerial();
    listenServer();
}


//接受串口数据发送到服务器
void listenSerial() {

    while (Serial.available() > 0) {
        comdata += char(Serial.read());
    }
    if (comdata != "")//如果接受到数据
    {
        client.print(comdata);//向服务器发送数据 
    }
    comdata = "";//清空数据
}


//连接服务器,接收信息,发送给串口
void listenServer() {
    //连接上服务器则led灯亮
    if (client.connected())
        digitalWrite(led, LOW);
    else
        digitalWrite(led, HIGH);

    //若未连接到服务端，则客户端进行连接。
    while (!client.connected()) {
        if (!client.connect(host, tcpPort)) {
            Serial.println("连接中....");
            delay(500);
        }
    }


    while (client.available()) {
//    char val = client.read();
        while (client.available() > 0) {
            val += char(client.read());
        }
        Serial.println(val);
        val = "";//清空数据
    }
    delay(1);
}
