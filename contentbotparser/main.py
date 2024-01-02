import asyncio
import logging
import os

from dotenv import load_dotenv

from clients.backend import APIClient
from parsers.parser import Parser


async def main():
    load_dotenv()
    logging.basicConfig(level="INFO")

    api_client = APIClient(os.getenv("API_SERVER"))
    parser = Parser(api_client)

    await asyncio.sleep(5)
    await parser.run()


if __name__ == "__main__":
    asyncio.run(main())
