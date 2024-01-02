import asyncio
import datetime

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient


async def pars_retailru(session: aiohttp.ClientSession, api: APIClient, timeout: int):
    url = "https://www.retail.ru/"
    while True:
        links = list()
        async with session.get(url) as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.find("div", class_="list-news")

            day = None
            for i in info.find_all(recursive=False):
                try:
                    if i.name == "p":
                        d = i.find("b").text
                        if d == "Сегодня":
                            day = datetime.datetime.now().day
                        elif d == "Вчера":
                            day = (datetime.datetime.now() - datetime.timedelta(days=1)).day
                        else:
                            break

                    if i.name == "ul":
                        for a in i.find_all(recursive=False):
                            t = datetime.datetime.strptime(
                                a.find("div", class_="date").text,
                                "%H:%M"
                            )

                            now = datetime.datetime.now()
                            content_datetime = datetime.datetime(
                                year=now.year,
                                month=now.month,
                                day=day,
                                hour=t.hour,
                                minute=t.minute
                            )

                            if content_datetime + datetime.timedelta(seconds=timeout) > now:
                                links.append(a.find("a").attrs["href"])

                except Exception:
                    continue

        contents = list()
        for link in links:
            try:
                async with session.get(url[:-1] + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.find("h1").text
                    info_body = soup.select("div.detail-articles")[0].find_all("p")

                    for i in info_body:
                        i = i.find("img")
                        if i is not None:
                            print(i.attrs["src"])
            except Exception:
                continue
        await asyncio.sleep(timeout)
