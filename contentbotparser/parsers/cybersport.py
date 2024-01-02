import asyncio
import datetime

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient
from utils import text_validator


async def pars_cybersportru(session: aiohttp.ClientSession, api: APIClient, timeout: int):
    url = "https://www.cybersport.ru"
    while True:
        links_imgs = dict()
        async with session.get(url + "/?sort=-publishedAt") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select("div.root_d51Rr")
            now = datetime.datetime.now()
            for i in info:
                try:
                    dt = datetime.datetime.strptime(i.find("span").text, "%d.%m в %H:%M")
                    content_datetime = datetime.datetime(
                        year=now.year,
                        month=dt.month,
                        day=dt.day,
                        hour=dt.hour,
                        minute=dt.minute
                    )

                    if content_datetime + datetime.timedelta(seconds=timeout) > now:
                        link = i.find("a").attrs["href"]
                        if i.find("img"):
                            links_imgs[link] = i.find("img").attrs["src"]
                        else:
                            links_imgs[link] = None

                except Exception:
                    continue

        contents = list()
        for link, img in links_imgs.items():
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.find("h1").text
                    info = soup.select(".text-content")[0]

                    txt = info_head + "\n\n"
                    for i in info.find_all("p"):
                        if i.attrs["class"] == ['authoredQuoteContent_U+5O3']:
                            txt += i.parent.find("span").text + ": \""
                            txt += i.text + "\"\n\n"
                        else:
                            txt += i.text + "\n\n"

                    if text_validator(txt, bool(img)):
                        data = {
                            "text": txt,
                            "img": [img],
                            "types": ["all", "cybersport"]
                        }
                        contents.append(data)

            except Exception:
                continue

        await api.send_content(contents)
        await asyncio.sleep(timeout)