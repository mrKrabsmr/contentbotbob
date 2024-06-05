import asyncio
import datetime
import logging

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient
from utils import text_validator


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
                            day = (datetime.datetime.now() -
                                   datetime.timedelta(days=1)).day
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

                except Exception as e:
                    logging.info(e)

        contents = list()
        for link in links:
            try:
                async with session.get(url[:-1] + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.find("h1").text
                    info_body = soup.select(
                        "div.detail-articles")[0].find_all("p")
                    info_body = info_body[:len(info_body)//2-1]

                    plain_txt = ""
                    img = None
                    for i in info_body:
                        im = i.find("img")
                        if im is not None and img is None:
                            src = im.attrs["src"] 
                            if not src.endswith("gif"):
                                img = url + im.attrs["src"]

                        el = getattr(i, "text", "")
                        plain_txt += el
                        if el:
                            plain_txt += "\n\n"

                    txt = ""
                    for char in plain_txt:
                        if char == "\n":
                            if txt[-2:] != "\n\n":
                                txt += char
                        else:
                            txt += char

                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": [img],
                            "types": ["all", "retails"]
                        }
                        contents.append(data)

            except Exception as e:
                logging.error(e)

        await api.send_content(contents)
        await asyncio.sleep(timeout)
