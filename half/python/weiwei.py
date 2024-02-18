from appium import webdriver

# 设置capabilities
capabilities = {
    'platformName': 'Android',
    'automationName': 'uiautomator2',
    'deviceName': '332b8b4d',
    'appPackage': 'com.dyheart.chat',
    'appActivity': 'com.dyheart.chat.MainActivity',
    'noReset': True
}

# 使用webdriver.Remote
driver = webdriver.Remote('http://localhost:4723/wd/hub', capabilities)
def send_emoji():
    # 打开聊天窗口，这里假设聊天窗口的元素 id 为 'chat_window'
    chat_window = driver.find_element_by_id('chat_windowwarp)
    chat_window.click()

    # 打开表情面板，这里假设表情面板的元素 id 为 'emoji_panel'
    emoji_panel = driver.find_element_by_id('emoji_panel')
    emoji_panel.click()

    # 选择并发送表情，这里假设表情的元素 id 为 'emoji'
    emoji = driver.find_element_by_id('emoji')
    emoji.click()

    # 点击发送按钮，这里假设发送按钮的元素 id 为 'send_button'
    send_button = driver.find_element_by_id('send_button')
    send_button.click()

# 每隔一段时间发送一次表情
while True:
    send_emoji()
    time.sleep(5)  # 等待一分钟