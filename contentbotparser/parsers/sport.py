import asyncio
import datetime

import aiohttp
from bs4 import BeautifulSoup

from clients.backend import APIClient
from utils import text_validator, months


async def pars_championatcom(session: aiohttp.ClientSession, api: APIClient, kind: str, timeout: int):
    url = "https://www.championat.com/"
    while True:
        links = list()
        async with session.get(url + f"news/{kind}/1.html") as resp:
            html = await resp.text()
            soup = BeautifulSoup(html, "html.parser")
            info = soup.select("div.news-items")[0]
            now = datetime.datetime.now()
            y, m, d = now.year, now.month, now.day
            for x in info.find_all("div"):
                try:
                    if x.attrs["class"][0] == "news-items__head":
                        val = x.text
                        for a, b in months.items():
                            val = val.replace(a, b)

                        ymd = datetime.datetime.strptime(val, "%d %B %Y")
                        y, m, d = ymd.year, ymd.month, ymd.day

                    if x.attrs["class"][0] == "news-item__time":
                        content_time = datetime.datetime.strptime(x.text, "%H:%M")

                        content_datetime = datetime.datetime(
                            y, m, d, content_time.hour, content_time.minute
                        )

                        if content_datetime + datetime.timedelta(seconds=timeout) > datetime.datetime.now():
                            links.append(x.parent.find("div", class_="news-item__content").find("a").attrs["href"])
                except Exception:
                    continue

        contents = list()
        for link in links:
            try:
                async with session.get(url + link) as resp:
                    html = await resp.text()
                    soup = BeautifulSoup(html, "html.parser")
                    info_head = soup.select(".article-head__title")[0]
                    info_body = soup.select(".article-content")[0]
                    info_img = soup.select(".article-head__photo")

                    txt = info_head.text + "\n\n"
                    for i in info_body.find_all("p"):
                        if "фото:" in i.lower():
                            continue
                        
                        txt += i.text + "\n\n"
                    
                    if txt[-1] == ":":
                        txt = "\n".join(txt.split("\n")[:-1])

                    img = None
                    if len(info_img) != 0:
                        img = info_img[0].find("img").attrs["src"]

                    if text_validator(txt, img is not None):
                        data = {
                            "text": txt,
                            "img": [img],
                            "types": ["sport", "all", kind]
                        }

                        contents.append(data)

            except Exception:
                continue

        await api.send_content(contents)
        await asyncio.sleep(timeout)

