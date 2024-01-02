"""
Описание моделей, связанных с пользователями
@version 1.0
"""
from django.contrib.auth.base_user import BaseUserManager
from django.contrib.auth.hashers import make_password
from django.contrib.auth.models import AbstractUser
from django.db import models

from apps.models import ModelCore
from apps.users.utils import get_activation_code


class UserManager(BaseUserManager):
    def create_user(self, username, email, password):
        if not username:
            raise ValueError("User must have a login")

        if not password:
            raise ValueError("User must have a password")

        user = self.model(username=username, email=email, password=password)
        user.save()

        return user

    def create_superuser(self, username, email, password):
        if not username:
            raise ValueError("User must have a username")

        if not password:
            raise ValueError("User must have a password")

        user = self.model(
            username=username, email=email, is_superuser=True, is_staff=True, password=password
        )
        user.save()

        return user


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

    objects = UserManager()

    def save(self, *args, **kwargs):
        if self._state.adding or self.password != self.__class__.objects.get(pk=self.pk).password:
            self.password = make_password(self.password)
        return super().save(*args, **kwargs)


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

