import asyncio

import aiohttp

from parsers.animals import pars_gismeteoru
from parsers.crypto import pars_investcom
from parsers.films import pars_filmru
from parsers.sport import pars_championatcom


class Parser:
    def __init__(self, api):
        self.api = api

    async def run(self):
        async with aiohttp.ClientSession() as session:
            session.headers.add(
                key="User-Agent",
                value="Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.3"
            )

            await asyncio.gather(
                pars_filmru(session=session, api=self.api, timeout=60 * 60 * 24),
                pars_gismeteoru(session=session, api=self.api, timeout=60*60*5),
                pars_investcom(session=session, api=self.api, timeout=60*60),
                pars_championatcom(session=session, api=self.api, kind="football", timeout=60 * 30),
                pars_championatcom(session=session, api=self.api, kind="hockey", timeout=60 * 60 * 2),
                pars_championatcom(session=session, api=self.api, kind="tennis", timeout=60 * 60 * 3),
                pars_championatcom(session=session, api=self.api, kind="boxing", timeout=60 * 60 * 2),
                pars_championatcom(session=session, api=self.api, kind="basketball", timeout=60 * 60 * 2),
                pars_championatcom(session=session, api=self.api, kind="auto", timeout=60 * 60),
                pars_championatcom(session=session, api=self.api, kind="volleyball", timeout=60 * 60 * 24),
                pars_championatcom(session=session, api=self.api, kind="figureskating", timeout=60 * 60 * 3),
            )
