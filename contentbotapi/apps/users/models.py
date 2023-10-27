"""
Описание моделей, связанных с пользователями
@version 1.0
"""

from django.contrib.auth.models import AbstractUser
from django.db import models

from apps.models import ModelCore
from apps.users.utils import get_activation_code


class User(ModelCore, AbstractUser):
    """
    Extend Django AbstractUser

    Deleted fields:
    - `first_name`
    - `last_name`
    Note: Enough to use `username`

    Added fields:
    - `is_confirmed` boolean field. Shows is `email` confirmed or not
    """

    is_confirmed = models.BooleanField(default=False)

    first_name = None
    last_name = None


class UserCodeActivation(ModelCore):
    """
    fields:
    - `user` foreign key field. The user with whom the activation code is associated
    - `code` string field. This is a six digit code sent to the user's email
    """
    user = models.OneToOneField(to=User, on_delete=models.CASCADE, related_name="code")
    code = models.CharField(max_length=255, default=get_activation_code)

    class Meta:
        db_table = "code_activations"
        verbose_name = "Код активации"
        verbose_name_plural = "Коды активации"

    def save(*args, **kwargs):
        pass
