from apps.channels.models import Channel
from apps.subscribes.services import UserSubscribeService as US_s


class ChannelService:
    _queryset = Channel.objects.all()

    @classmethod
    def get_user_channels(cls, user):
        return cls._queryset.filter(owner=user)

    @classmethod
    def get_filtered_channels_for_contents(cls, **kwargs):
        channels = cls._queryset.filter(
            settings__allowed_content_sources__value=kwargs.get("source"),
            settings__min_rating__lte=kwargs.get("rating")
        )

        return channels

    @classmethod
    def check_user_may_add_channel(cls, user):
        user_subscribe = US_s.get_user_subscribe(user)
        if not user_subscribe:
            return False

        subscribe = user_subscribe.subscribe
        user_channels_count = cls._queryset.filter(owner=user).count()

        if subscribe.max_use_channels <= user_channels_count:
            return False

        return True
