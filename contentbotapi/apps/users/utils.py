from random import sample
from string import digits


def send_email_after_registration():
    pass


def get_activation_code():
    return "".join(sample(digits, 6))
