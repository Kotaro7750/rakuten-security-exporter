from selenium import webdriver
from selenium.webdriver.common.by import By
import time
import os


def setup_driver(download_dir):
    options = webdriver.ChromeOptions()
    options.add_argument('--headless=new')
    options.add_argument('--no-sandbox')
    options.add_argument('--disable-dev-shm-usage')
    options.add_experimental_option('prefs', {
        "download.default_directory": download_dir
    })
    driver = webdriver.Chrome(options)
    driver.set_window_size(1920, 1080)
    driver.implicitly_wait(5)

    return driver


def login_and_expand_mymenu(driver, id, password):
    driver.get('https://www.rakuten-sec.co.jp/')

    # ログイン
    driver.find_element(By.ID, 'form-login-id').send_keys(id)
    driver.find_element(By.ID, 'form-login-pass').send_keys(password)
    driver.find_element(By.ID, 'login-btn').click()

    # マイメニュー展開
    driver.find_element(By.XPATH, '/html/body/header/div/div[4]/div/nav/ul/li[6]/button').click()


def download_withdrawal_history(id, password, download_dir):
    print("Start downloading withdrawal history")
    print("Download dir is {}".format(download_dir))
    driver = setup_driver(download_dir)
    login_and_expand_mymenu(driver, id, password)

    driver.find_element(By.LINK_TEXT, '入出金履歴').click()
    driver.find_element(By.XPATH, '//img[@alt="CSVで保存"]').click()

    time.sleep(5)
    driver.quit()
    print("Finish downloading withdrawal history")


def download_dividened_history(id, password, download_dir):
    print("Start downloading dividend history")
    print("Download dir is {}".format(download_dir))
    driver = setup_driver(download_dir)
    login_and_expand_mymenu(driver, id, password)

    driver.find_element(By.LINK_TEXT, '配当・分配金').click()
    driver.find_element(By.XPATH, '//img[@alt="すべて"]').click()

    driver.find_element(By.XPATH, '/html/body/div[2]/div/div/div/div/table/tbody/tr/td/form/div[3]/table/tbody/tr[7]/td/span/input').click()
    driver.find_element(By.XPATH, '//img[@alt="CSVで保存"]').click()

    time.sleep(5)
    driver.quit()
    print("Finish downloading dividend history")


def download_asset_list(id, password, download_dir):
    print("Start downloading asset list")
    print("Download dir is {}".format(download_dir))
    driver = setup_driver(download_dir)

    try:
        login_and_expand_mymenu(driver, id, password)

        time.sleep(1)
        driver.find_element(By.LINK_TEXT, '保有商品一覧').click()
        # 読み込みが終わってからCSVで保存する
        time.sleep(5)
        driver.find_element(By.XPATH, '//img[@alt="CSVで保存"]').click()

        time.sleep(5)

    except:
        driver.get_screenshot_as_file('hoge')

    driver.quit()

    print("Finish downloading asset list")


id = os.environ['RAKUTEN_SEC_ID']
password = os.environ['RAKUTEN_SEC_PASSWORD']
download_dir = os.environ['DOWNLOAD_DIR']

download_withdrawal_history(id, password, download_dir)
download_dividened_history(id, password, download_dir)
download_asset_list(id, password, download_dir)
