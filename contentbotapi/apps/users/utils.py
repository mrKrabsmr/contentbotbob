from random import sample
from string import digits

from apps.users.models import User, UserCodeActivation


def send_email_after_registration():
    pass


def create_or_change_activation_code(user):
    code_activation, created = UserCodeActivation.objects.get_or_create(user=user)
    if not created:
        code_activation.code = "".join(sample(digits, 6))
        code_activation.save()
