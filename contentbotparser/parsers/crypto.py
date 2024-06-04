import asyncio
import logging

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
            info = soup.select("#__next > div.md\:relative.md\:bg-white > div.relative.flex > div.grid.flex-1.grid-cols-1.px-4.pt-5.font-sans-v2.text-\[\#232526\].antialiased.xl\:container.sm\:px-6.md\:grid-cols-\[1fr_72px\].md\:gap-6.md\:px-7.md\:pt-10.md2\:grid-cols-\[1fr_420px\].md2\:gap-8.md2\:px-8.xl\:mx-auto.xl\:gap-10.xl\:px-10 > div.min-w-0 > div > div.news-analysis-v2_articles-container__3fFL8.mdMax\:px-3.mb-12 > ul > li:nth-child(n) > article > div > ul > li > div > time")
            for i in info:
                try:
                    if "мин" in i.text or "сек" in i.text:
                        link = i.parent.parent.parent.parent.find("a").attrs["href"]
                        links.append(link)

                except Exception:
                    continue
        contents = list()
        for link in links:
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info = soup.select("#article > div")[0]

                    plain_txt = ""
                    for i in info.find_all("p"):
                        try:
                            plain_txt += i.text + "\n\n"
                        except Exception:
                            continue

                    txt = soup.find("h1").text + "\n\n" + plain_txt.replace(
                        "Happycoin.club - ", ""
                    ).split(
                        "Позиция успешно добавлена:"
                    )[0].split(
                        "Источник: "
                    )[0].split(
                        "Читайте оригинальную статью на сайте Happycoin.club"
                    )[0].strip()

                    img = soup.select("#__next > div.md\:relative.md\:bg-white > div.relative.flex > div.grid.flex-1.grid-cols-1.px-4.pt-5.font-sans-v2.text-\[\#232526\].antialiased.xl\:container.sm\:px-6.md\:grid-cols-\[1fr_72px\].md\:gap-6.md\:px-7.md\:pt-10.md2\:grid-cols-\[1fr_420px\].md2\:gap-8.md2\:px-8.xl\:mx-auto.xl\:gap-10.xl\:px-10 > div.min-w-0 > div > div:nth-child(1) > div.relative.flex.flex-col > div.mb-5.mt-4.sm\:mt-8.md\:mb-8.relative.h-\[294px\].w-full.overflow-hidden.bg-\[\#181C21\].sm\:h-\[420px\].xl\:h-\[441px\] > img")[0].attrs["src"]
                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": [img],
                            "types": ["all", "crypto"]
                        }

                        contents.append(data)

            except Exception as e:
                logging.info(e)
        logging.info(contents)
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
