from app.loader import config


def build_endpoint_v1(route: str) -> str:
    return config.api_server + "api/v1/" + route
