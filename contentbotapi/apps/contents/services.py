from django.db.transaction import atomic

from apps.channels.services import ChannelService as Ch_s
from apps.contents.models import Content, ContentMedia


class ContentService:
    _queryset = Content.objects.all()

    @classmethod
    def get_suitable_one(cls, outer_id):
        suitable_content = cls._queryset.select_related(
            "channel"
        ).prefetch_related(
            "images"
        ).filter(
            channel__outer_id=outer_id
        ).order_by(
            "created_at"
        ).first()

        return suitable_content


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
            content = Content(text=data["text"], rating=data["rating"], channel=channel)
            content_image = ContentMedia(file_url=data["img"], channel_content=content)

            contents.append(content)
            content_images.append(content_image)

        cls._queryset.bulk_create(contents)
        cls._media_queryset.bulk_create(content_images)

    @staticmethod
    def _filter_channels(data):
        empty = False if data.get("img") else True

        channels = Ch_s.get_filtered_channels_for_contents(source=data.get("source"), rating=data.get("rating"))

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
