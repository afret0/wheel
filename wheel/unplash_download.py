import requests
import time, os


class UnspalshDownload():
    def __init__(self):
        self.host = 'https://unsplash.com'

    def get_imgs_list(self, page):
        target = '{}/napi/collections/1065976/photos'.format(self.host)
        headers = {}
        params = {'page': page, 'per_page': 10, 'order_by': 'latest'}
        try:
            r = requests.get(url=target, params=params)
        except Exception as e:
            print('获取 img{} 列表失败 ... {}'.format(page, e))
            # raise e
        else:
            # print(r.json())
            imgs = [{'id': i['id'], 'download': i['links']['download']} for i in r.json()]
            print('get img list ...')
            return imgs

    def save_img(self, img, filepath):
        IMG = '{}.jpg'.format(img['id'])
        print('saving  {} ...'.format(IMG))
        if not os.path.exists(filepath):
            os.mkdir(filepath)
        if os.path.exists(IMG):
            try:
                r = requests.get(url=img['download'])
            except Exception as e:
                print('下载 {} 失败 ... {}'.format(IMG, e))
                # raise e
            else:
                if not os.path.exists(IMG):
                    with open(os.path.join(filepath, IMG), 'wb') as f:
                        f.write(r.content)
            print('saved  {} ...'.format(IMG))
        else:
            print('{} 已下载 ...'.format(IMG))


def start(page):
    USR_PATH = os.environ.get('USERPROFILE')
    FIEL_PATH = os.path.join(os.path.join(USR_PATH, 'Downloads'), 'unsplash')
    d = UnspalshDownload()
    IMGS = d.get_imgs_list(page)
    for IMG in IMGS:
        d.save_img(IMG, FIEL_PATH)


if __name__ == '__main__':
    from multiprocessing import Pool

    START_PAGE = 10
    END_PAGE = 100
    pool = Pool()
    page = [i for i in range(START_PAGE, END_PAGE + 1)]

    pool.map(start, page)
    pool.close()
    pool.join()

    pass
