import asyncio
import datetime

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient
from utils import text_validator, months


async def pars_filmru(session: aiohttp.ClientSession, api: APIClient, timeout: int):
    url = "https://www.film.ru/"
    while True:
        links = list()
        async with session.get(url + "topic/news") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select("div.redesign_topic > div:nth-child(2)")
            for i in info:
                try:
                    now = datetime.datetime.now()
                    dt = i.find_all(text=True, recursive=False)[-1].strip()

                    if "сегодня" in dt.lower():
                        d = now.day
                        hm = datetime.datetime.strptime(dt, "Сегодня %H:%M")
                        d, h, m = d, hm.hour, hm.minute
                    else:
                        for x, y in months.items():
                            dt = dt.replace(x, y)

                        dhm = datetime.datetime.strptime(dt, "%d %B %H:%M")
                        d, h, m = dhm.day, dhm.hour, dhm.minute

                    content_datetime = datetime.datetime(
                        year=now.year,
                        month=now.month,
                        day=d,
                        hour=h,
                        minute=m,
                    )

                    if content_datetime + datetime.timedelta(seconds=timeout) > now:
                        links.append(i.parent.attrs["onclick"].split("'")[1])

                except Exception:
                    continue

        contents = list()
        for link in links:
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.select(".wrapper_articles_left > h1:nth-child(2)")[0]
                    img = url + info_head.parent.find("div", class_="wrapper_articles_background").find("img").attrs[
                        "src"]

                    info_body = soup.select(".wrapper_articles_text")[0]
                    plain_txt = ""
                    for i in info_body.find_all("p"):
                        plain_txt += i.text

                    txt = info_head.text + "\n\n" + plain_txt.split("Источник: ")[0]

                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": [img] if img else None,
                            "types": ["all", "films"]
                        }

                        contents.append(data)

            except Exception:
                continue

        await api.send_content(contents)
        await asyncio.sleep(timeout)
