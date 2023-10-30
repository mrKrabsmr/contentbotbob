import uuid

from django.db import models

from apps.channels.models import Channel


class Content(models.Model):
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
        verbose_name = "Контент"
        verbose_name_plural = "Контенты"
        ordering = ("-created_at", "rating")

    def __str__(self):
        return self.text[:25] + "..."


class ContentMedia(models.Model):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    channel_content = models.ForeignKey(
        to=Content, on_delete=models.CASCADE, related_name="images"
    )
    file_url = models.CharField(max_length=255, null=True, blank=True)

    class Meta:
        db_table = "channel_content_files"
        verbose_name = "Контент-Медиа"
        verbose_name_plural = "Контент-Медиа"

    def __str__(self):
        return str(self.channel_content)
