import os

from aiogram.contrib.fsm_storage.memory import MemoryStorage
from aiogram.types import ParseMode


class Config:
    def __init__(self):
        self.api_server = os.getenv("API_SERVER")
        self.bot_token = os.getenv("API_TOKEN")
        self.parse_mode = ParseMode.HTML
        self.storage = MemoryStorage
