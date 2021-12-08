#include "RC522.h"
#include <SPI.h>
#include <EEPROM.h>
#include<Servo.h>
#include <ArduinoJson.h>
#include <SoftwareSerial.h>
#include <LiquidCrystal_I2C.h>

#define N 4

Servo myservo;
StaticJsonDocument<256> doc;
SoftwareSerial espSerial(2, 3); // RX, TX
LiquidCrystal_I2C lcd(0x27, 16, 2);

int addr = 2 * N;
String userName[N];
String cardId[N];
String strA;
int mutex = 1;

void show(){
   for(int i=0;i<4;i++){
      Serial.println(userName[i]);
      Serial.println(cardId[i]);
    }
}
void setup() {
    //每次启动 先读取EEPROM里的数据更新用户
     Serial.begin(9600);
    espSerial.begin(9600);
    Serial.println("begin...");
    lcd.init(); 
    lcd.backlight();
    mutex = 1;
   
    SPI.begin();
    pinMode(chipSelectPin, OUTPUT);
    digitalWrite(chipSelectPin, LOW);
    pinMode(NRSTPD, OUTPUT);
    MFRC522_Init();

    myservo.attach(7);
    myservo.write(0);

    pinMode(7, OUTPUT);
    pinMode(4, OUTPUT);
}


void loop() {
    myRead();
    if (mutex == 1) {
        readCard();
    }
    receive();

}


void receive() {
    if (espSerial.available() > 0) {
        mutex = 0;
        if (espSerial.peek() != '*') {
            strA += (char) espSerial.read();
        } else {
            espSerial.flush();
            onlineOpen();
            while (espSerial.available() > 0) espSerial.read();
            strA = "";
            mutex = 1;
        }
    }
}


//网络接收消息
void onlineOpen() {
    Serial.println(strA);
    //接收到"open"开门
    if (strA.indexOf("open") != -1) {
        openDoor();
        //lcd使用的内存空间可能被占用了
        //不重新init就一直乱码
        //只有每次使用lcd显示就init一次
        lcd.init();
        lcd.clear();
        lcd.setCursor(4, 0);
        lcd.print("welcome");
        lcd.setCursor(0, 1);
        lcd.print("the door is open");
        delay(30);
    } else {
        //解析json数据
        //创建DynamicJsonDocument对象
        DeserializationError error = deserializeJson(doc, strA);


        if (error) {
//            Serial.print(F("deserializeJson() failed: "));
//            Serial.println(error.c_str());
            return;
        }

        int count = 0;
        for (JsonObject item : doc.as<JsonArray>()) {

            const char *user_name = item["user_name"];
            const char *card_id = item["card_id"];
            userName[count] = user_name;
            cardId[count] = card_id;
            userName[count].trim();
            cardId[count].trim();
//            Serial.print(userName[count]);
//            Serial.println(cardId[count]);
            count++;
        }
        //把更新的数据保存到eeprom中
        myWrite();
        show();
        myRead();
        show();
    }


}

//读卡器开门
void readCard() {
    char serNum[5];
    unsigned char status;
    unsigned char str[10];
    status = MFRC522_Request(PICC_REQIDL, str);
    if (status == MI_OK)      //读取到ID卡时候
    {
        status = MFRC522_Anticoll(str);
        if (status == MI_OK) {
            memcpy(serNum, str, 5);
            String card = toString(serNum);
            int i;
            for (i = 0; i < 4; i++) {
                if (cardId[i].compareTo(card) == 0) {
                    openDoor();
                    String jsonData = jsonMarshal(userName[i], cardId[i]);
                    espSerial.println(jsonData);
                    
                    lcd.init();
                    showCardID(i);
                    delay(30);
                    
                    break;
                }
            }
            if (i == 4) {
//                 Serial.println(card);
//                 Serial.println("非法用户!");

                lcd.init();
                
                lcd.setCursor(0, 0);
                lcd.print(card);

                lcd.setCursor(0, 1);
                lcd.print("invalid card");
            }

        }
        delay(1300);
    }

    MFRC522_Halt();

}

//开门
void openDoor() {
    myservo.write(270);
//    Serial.println("The door has been opened!");
//    Serial.println();
    delay(1500);
    myservo.write(0);
}

void showCardID(int id) {
    lcd.setCursor(0, 0);
    lcd.print(userName[id]);
    lcd.setCursor(0, 1);
    lcd.print(cardId[id]);

}

//把读取到的卡片转为十六进制的字符串
String toString(char id[5]) {
    char _16[] = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'};
    int IDlen = 4;
    String str;
    for (int i = 0; i < IDlen; i++) {
        str += tohex((0x0F & id[i] >> 4));
        str += tohex((0x0F & id[i]));
    }
    return str;
}


//十进制转十六进制
String tohex(int n) {
    if (n == 0) return "0";
    String result = "";
    char _16[] = {
            '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'
    };
    const int radix = 16;
    while (n) {
        int i = n % radix;
        result = _16[i] + result;
        n /= radix;
    }

    return result;
}


//写入EEPROM
void myWrite() {
    //EEPROM.put写入后断电重启不知道为什么会乱码
    //所以用write 把每个字符转为ascil后写入
    int addr = 2*N;
    //username 
    for (int i = 0; i < N; i++) {
        int count = 0;
        for (char temp:userName[i]) {
            int x = toascii(temp);
            EEPROM.write(addr++, x);
            count++;
        }
        EEPROM.write(i, count);
    }

    //cardId
    for (int i = 0; i < N; i++) {
        int count = 0;
        for (char temp:cardId[i]) {
            int x = toascii(temp);
            EEPROM.write(addr++, x);
            count++;
        }
        EEPROM.write(N + i, count);
    }
}

//读取EEPROM
void myRead() {
    //读取ascil再转换
    int addr = 2 * N;
    //username
    for (int i = 0; i < N; i++) {
        String strTemp;
        int len = EEPROM.read(i);
        for (int j = 0; j < len; j++) {
            strTemp += (char) EEPROM.read(addr++);
        }
        userName[i] = strTemp;
    }

    //cardId
    for (int i = 0; i < N; i++) {
        String strTemp;
        int len = EEPROM.read(N + i);
        for (int j = 0; j < len; j++) {
            strTemp += (char) EEPROM.read(addr++);
        }
        cardId[i] = strTemp;
    }
}

//包装成json发送
String jsonMarshal(String user_name, String card_id) {
    user_name = "\"" + user_name + "\"";
    String str = "{";
    str += "\"user_name\":";
    str += user_name;
    str += ",";
    str += "\"card_id\":\"";
    str += card_id;
    str += "\"}";
    return str;
}
