import uuid

from django.db import models

from apps.users.models import User


class Subscribe(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    name = models.CharField(max_length=255)
    description = models.TextField(null=True, blank=True)
    # period_in_days = models.PositiveIntegerField(default=30, null=True, blank=True)
    max_use_channels = models.PositiveSmallIntegerField()
    price_rub = models.DecimalField(default=0.0, decimal_places=2, max_digits=8)

    class Meta:
        db_table = "subscribes"
        verbose_name = "подписка"
        verbose_name_plural = "подписки"

    def __str__(self):
        return self.name


class UserSubscribe(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    user = models.OneToOneField(to=User, on_delete=models.CASCADE, related_name="subscribe")
    subscribe = models.ForeignKey(to=Subscribe, on_delete=models.SET_NULL, related_name="user", null=True)
    bought_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "users_subscribes"
        verbose_name = "подписка пользователя"
        verbose_name_plural = "подписки пользователей"
        ordering = ("-bought_at",)
