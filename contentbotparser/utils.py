months = {
    "января": "January",
    "февраля": "February",
    "марта": "March",
    "апреля": "April",
    "мая": "May",
    "июня": "June",
    "июля": "July",
    "августа": "August",
    "сентября": "September",
    "октября": "October",
    "ноября": "November",
    "декабря": "December",
}


def text_validator(txt: str, has_img: bool) -> bool:
    if has_img:
        if len(txt) > 1032:
            return False
    else:
        if len(txt) > 4096:
            return False

    return True
