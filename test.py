import requests
import random
import string
from faker import Faker

BASE_URL = "http://localhost:4220"





# 生成随机英文字符串
def random_english_string(length):
    # 选择所有字母和数字作为字符集
    characters = string.ascii_letters + string.digits
    # 使用random.choices()函数生成指定长度的随机字符串
    return ''.join(random.choices(characters, k=length))

# 生成随机中文字符串
def random_chinese_string(length):
    result = ''
    for i in range(length):
        # 随机选择Unicode编码值，转换为中文字符
        result += chr(random.randint(0x4e00, 0x9fbf))
    return result


def random_integer():
    # 生成1~5位的随机整数
    return random.randint(1, 10**5)

def test_add_song(name, user, uid):
    url = BASE_URL + "/songs/add"
    payload = {"name": name, "user": user, "uid": uid}
    response = requests.post(url, data=payload)
    print("add_song: ", response.status_code, response.text)

def test_list_songs():
    url = BASE_URL + "/songs"
    response = requests.get(url)
    print("list_songs: ", response.status_code, response.text)

def test_delete_song(number):
    url = BASE_URL + f"/songs/del"
    payload = {"number": number}
    response = requests.post(url, data=payload)
    print("delete_song: ", response.status_code, response.text)

def test_clear_song():
    url = BASE_URL + f"/songs/clear"
    payload = {}
    response = requests.post(url, data=payload)
    print("clear_song: ", response.status_code, response.text)



def add_test_songs(n):
    for i in range(n):
        f1 = Faker(['en_US'])
        f2 = Faker(['zh_CN'])

        test_add_song(f1.name(),f2.name(),random_integer())


if __name__ == "__main__":
    add_test_songs(5)