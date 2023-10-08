import uuid

from django.conf import settings
from django.db import models

from apps.users.models import User


class Channel(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    owner = models.ForeignKey(to=User, on_delete=models.CASCADE, related_name="channels")
    name = models.CharField(max_length=255)
    type = models.CharField(
        max_length=255, choices=[(i, i) for i in settings.CHANNEL_TYPE_LIST]
    )
    resource = models.CharField(
        max_length=255, choices=[(i, i) for i in settings.RESOURCE_LIST]
    )
    outer_id = models.CharField(max_length=255)

    class Meta:
        db_table = "channels"
        verbose_name = "Канал"
        verbose_name_plural = "Каналы"

    def __str__(self):
        return str(self.name)


class ChannelSettings(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    channel = models.ForeignKey(
        to=Channel, on_delete=models.CASCADE, related_name="settings"
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


class ChannelSettingsAllowedContentSource(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
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


class ChannelContent(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    channel = models.ForeignKey(
        to=Channel, on_delete=models.CASCADE, related_name="contents"
    )
    text = models.TextField(null=True, blank=True)
    rating = models.IntegerField(
        null=True, blank=True, choices=[(i, i) for i in range(1, 11)]
    )
    created_at = models.DateTimeField(auto_now_add=True)

    class Meta:
        db_table = "channel_contents"
        verbose_name = "Каналы-Контент"
        verbose_name_plural = "Каналы-Контенты"
        ordering = ("-created_at", "rating")

    def __str__(self):
        return self.text[:25] + "..."


class ChannelContentMedia(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    channel_content = models.ForeignKey(
        to=ChannelContent, on_delete=models.CASCADE, related_name="images"
    )
    file_url = models.CharField(max_length=255, null=True, blank=True)

    class Meta:
        db_table = "channel_content_files"
        verbose_name = "Каналы-Контент-Медиа"
        verbose_name_plural = "Каналы-Контент-Медиа"

    def __str__(self):
        return str(self.channel_content)
