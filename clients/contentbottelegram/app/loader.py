import aiogram.contrib.fsm_storage.memory
from aiogram import Bot, Dispatcher
from dotenv import load_dotenv

from config import Config

load_dotenv()

config = Config()

bot = Bot(token=config.bot_token, parse_mode=config.parse_mode)

dp = Dispatcher(bot=bot, storage=config.storage)
