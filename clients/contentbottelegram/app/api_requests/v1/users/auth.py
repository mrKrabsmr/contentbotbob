import aiohttp

from app.api_requests.v1 import build_endpoint_v1


async def user_register(data):
    endpoint = build_endpoint_v1("register/")

    async with aiohttp.ClientSession() as session:
        response = await session.post(endpoint, data=data)
        response_text = response.json()
        if response.status != 201:
            return False, response_text

        return True, response_text


async def user_authenticate(data):
    pass
