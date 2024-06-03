import asyncio

import aiohttp
from bs4 import BeautifulSoup
from bs4.element import ResultSet

from clients.backend import APIClient
from utils import text_validator


async def pars_investcom(session: aiohttp.ClientSession, api: APIClient, timeout: int):
    url = "https://ru.investing.com/"
    while True:
        links = list()
        async with session.get(url + "news/cryptocurrency-news") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select(".largeTitle > article > div > span > span:nth-of-type(2)")
            for i in info:
                try:
                    if "минут" in i.text or "секунд" in i.text:
                        link = i.parent.parent.find("a").attrs["href"]
                        links.append(link)

                except Exception:
                    continue

        contents = list()
        for link in links:
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info = soup.select("div.WYSIWYG:nth-child(8)")[0]

                    plain_txt = ""
                    for i in info.find_all("p"):
                        plain_txt += i.text + "\n\n"

                    txt = info.parent.find("h1").text + "\n\n" + plain_txt.replace(
                        "Happycoin.club - ", ""
                    ).split(
                        "Позиция успешно добавлена:"
                    )[0].split(
                        "Источник: "
                    )[0].split(
                        "Читайте оригинальную статью на сайте Happycoin.club"
                    )[0].strip()

                    img = info.find("div", class_="imgCarousel").find("img").attrs["src"]

                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": [img],
                            "types": ["all", "crypto"]
                        }

                        contents.append(data)

            except Exception:
                continue

        await api.send_content(contents)
        await asyncio.sleep(timeout)


async def pars_finam(session, api: APIClient):
    url = "https://www.finam.ru/"
    while True:
        links = list()

        async with session.get(url + "publications/section/cryptonews/") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select("#finfin-local-plugin-block-item-publication-list-filter-date-content")
            print(info)
            for i in info:
                print(i.text)

        await asyncio.sleep(60 * 60 * 24)
