from aiogram import types
from aiogram.dispatcher.filters import CommandStart

from app.loader import dp


@dp.message_handler(CommandStart())
async def command_start(message: types.Message):
    pass