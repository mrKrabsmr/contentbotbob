import http
import logging

import aiohttp


class APIClient:
    def __init__(self, address):
        self.session = aiohttp.ClientSession()
        self.address = address

    async def send_content(self, data):
        url = self.address + "parser/contents/news/"
        async with self.session.post(url, json=data) as resp:
            if resp.status != http.HTTPStatus.CREATED:
                data = await resp.json()
                logging.error(data)



