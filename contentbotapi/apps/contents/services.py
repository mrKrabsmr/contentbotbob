import datetime
import logging

from django.db.transaction import atomic
from django.db.models import F

from apps.channels.services import ChannelService as Ch_s
from apps.channels.services import ChannelSettingsService as ChS_s
from apps.contents.models import Content, ContentMedia


class ContentService:
    _queryset = Content.objects.all()

    @classmethod
    def get_suitable_one(cls, outer_id) -> (Content, bool):
        channel = Ch_s.get_by_outer_id(outer_id)
        if not channel:
            return None, False

        if not channel.status_on:
            return None, False

        delta = datetime.datetime.now() - datetime.timedelta(days=ChS_s.get_by_channel_id(outer_id).not_later_days)
        cls._queryset.filter(
            channel__outer_id=outer_id, 
            created_at__lt=delta
        ).delete()

        suitable_content = cls._queryset.select_related(
            "channel"
        ).prefetch_related(
            "images"
        ).filter(
            channel__outer_id=outer_id,
            was_sent=False
        ).order_by(
            "-rating"
        ).first()

        if suitable_content:
            suitable_content.was_sent = True
            suitable_content.save()
            return suitable_content, True

        return None, True


class ContentDistributionService:
    _queryset = Content.objects.all()
    _media_queryset = ContentMedia.objects.all()

    @classmethod
    @atomic
    def saving_data(cls, data):
        data["rating"] = cls._evaluate_content(data)

        contents = []
        content_images = []

        for channel in cls._filter_channels(data):
            content = Content(text=data["text"],
                              rating=data["rating"], channel=channel)
            imgs = data["img"]
            if imgs:
                for img in imgs:
                    if img:
                        content_image = ContentMedia(
                            file_url=img, content=content)
                        content_images.append(content_image)

            contents.append(content)

        cls._queryset.bulk_create(contents)
        cls._media_queryset.bulk_create(content_images)

    @staticmethod
    def _filter_channels(data):
        empty = True
        imgs = data.get("img")
        if imgs and len(imgs) > 0:
            empty = False

        channels = Ch_s.get_filtered_channels_for_contents(
            type=data.get("types"),
            rating=data.get("rating"),
            text_part=data.get("text")[:50]
        )

        if empty:
            channels = channels.filter(settings__empty_file_allowed=True)

        return channels

    @staticmethod
    def _evaluate_content(data):
        rating = 0

        if data.get("img"):
            rating += 3

        len_text = len(data.get("text"))
        if len_text < 50 or len_text > 1500:
            rating += 1
        elif len_text < 100 or len_text > 1000:
            rating += 2
        elif len_text < 200 or len_text > 900:
            rating += 3
        elif len_text < 300 or len_text > 800:
            rating += 4
        elif len_text < 350 or len_text > 750:
            rating += 5
        elif len_text < 400 or len_text > 700:
            rating += 6
        else:
            rating += 7

        return rating
