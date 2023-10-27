from apps.channels.models import Channel
from apps.services import ServiceCore


class ChannelService(ServiceCore):
    _queryset = Channel.objects.all()

    @classmethod
    def get_user_channels(cls, user):
        return cls._queryset.filter(owner=user)

    @classmethod
    def get_filtered_channels_for_contents(cls, **kwargs):
        channels = cls._queryset.objects.filter(
            settings__allowed_content_sources__value=kwargs.get("source"),
            settings__min_rating__lte=kwargs.get("rating")
        )

        return channels
