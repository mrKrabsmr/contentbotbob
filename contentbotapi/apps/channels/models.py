import uuid

from django.conf import settings
from django.db import models

from apps.models import ModelCore
from apps.users.models import User


class Channel(ModelCore):
    owner = models.ForeignKey(to=User, on_delete=models.CASCADE, related_name="channels", blank=True)
    name = models.CharField(max_length=255)
    type = models.CharField(
        max_length=255, choices=[(i, i) for i in settings.CHANNEL_TYPE_LIST]
    )
    resource = models.CharField(
        max_length=255, choices=[(i, i) for i in settings.RESOURCE_LIST]
    )
    outer_id = models.CharField(max_length=255)
    status_on = models.BooleanField(default=True)

    class Meta:
        db_table = "channels"
        verbose_name = "Канал"
        verbose_name_plural = "Каналы"

    def __str__(self):
        return str(self.name)


class ChannelSettings(ModelCore):
    channel = models.OneToOneField(
        to=Channel, on_delete=models.CASCADE, related_name="settings", blank=True
    )
    min_rating = models.IntegerField(
        null=True, blank=True, choices=[(i, i) for i in range(1, 11)], default=1
    )
    empty_file_allowed = models.BooleanField(default=True)
    not_later_days = models.IntegerField(default=3)
    language = models.CharField(max_length=255, choices=[(i, i) for i in settings.CONTENT_LANGUAGES_LIST])

    class Meta:
        db_table = "channel_settings"
        verbose_name = "Настройки-Канал"
        verbose_name_plural = "Настройки-Каналы"

    def __str__(self):
        return str(self.channel)


class ChannelSettingsAllowedContentSource(ModelCore):
    channel_settings = models.ForeignKey(
        to=ChannelSettings, on_delete=models.CASCADE, related_name="allowed_content_sources"
    )
    value = models.CharField(
        max_length=255, choices=[(i, i) for i in settings.CONTENT_SOURCE_LIST]
    )

    class Meta:
        verbose_name = "Настройки канала - Разрешенный источник"
        verbose_name_plural = "Настройки каналов - Разрешенные источники"

    def __str__(self):
        return str(self.channel_settings)
