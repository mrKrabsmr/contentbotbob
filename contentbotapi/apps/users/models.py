import uuid
from random import sample
from string import digits

from django.contrib.auth.models import AbstractUser
from django.db import models


class User(AbstractUser):
    is_confirmed = models.BooleanField(default=False)

    first_name = None
    last_name = None


class UserCodeActivation(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.OneToOneField(to=User, on_delete=models.CASCADE, related_name="code")
    code = models.CharField(max_length=255, default="".join(sample(digits, 6)))

    class Meta:
        db_table = "code_activations"
        verbose_name = "Код активации"
        verbose_name_plural = "Коды активации"

    def save(*args, **kwargs):
        pass

