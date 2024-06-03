import asyncio
import datetime

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient
from utils import text_validator


async def pars_gismeteoru(session: aiohttp.ClientSession, api: APIClient, timeout: int):
    url = "https://www.gismeteo.ru/"
    while True:
        links = list()
        async with session.get(url + "news/animals/") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select("a.article-item")
            for i in info:
                try:
                    content_datetime = datetime.datetime.strptime(
                        i.attrs["data-pub-date"],
                        "%Y-%m-%dT%H:%M:%S"
                    )

                    if content_datetime + datetime.timedelta(seconds=timeout) > datetime.datetime.now():
                        links.append(i.attrs["href"])

                except Exception:
                    continue

        contents = list()
        for link in links:
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.select(".article-title > h1:nth-child(1)")[0]
                    info_body = soup.select(".article-content > div:nth-child(1)")[0]

                    txt = info_head.text + "\n\n"
                    for i in info_body.find_all("p"):
                        txt += i.text + "\n\n"

                    img = list()
                    for i in info_body.find_all("figure"):
                        img.append(i.find("img").attrs["src"])

                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": img if len(img) > 0 else None,
                            "types": ["all", "animals"]
                        }

                        contents.append(data)

            except Exception:
                continue

        await api.send_content(contents)
        await asyncio.sleep(timeout)
