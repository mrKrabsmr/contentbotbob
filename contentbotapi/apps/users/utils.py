from random import sample
from string import digits

from celery import shared_task
from django.core.mail import send_mail

from config.settings import EMAIL_HOST_USER, DEBUG


@shared_task
def send_code(email, code):
    return send_mail(
        subject="Код активации",
        message=code,
        from_email=EMAIL_HOST_USER,
        recipient_list=[email],
        fail_silently=DEBUG
    )


def get_activation_code():
    return "".join(sample(digits, 6))
