from apps.channels.models import Channel, ChannelSettings
from apps.subscribes.services import UserSubscribeService as US_s


class ChannelService:
    _queryset = Channel.objects.all()

    @classmethod
    def get_channels(cls, status_on=None):
        if status_on is None:
            return cls._queryset

        if isinstance(status_on, bool):
            return cls._queryset.filter(status_on=status_on)

        return None

    @classmethod
    def get_user_channels(cls, user):
        return cls._queryset.filter(owner=user)

    @classmethod
    def get_by_outer_id(cls, outer_id):
        return cls._queryset.filter(outer_id=outer_id).first()

    @classmethod
    def get_filtered_channels_for_contents(cls, **kwargs):
        channels = cls._queryset.filter(
            type__in=kwargs.get("type"),
            settings__min_rating__lte=kwargs.get("rating")
        ).exclude(
            contents__text__istartswith=kwargs.get("text_part")
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


class ChannelSettingsService:
    _queryset = ChannelSettings.objects.all()

    @classmethod
    def get_by_channel_id(cls, outer_id):
        channel = ChannelService.get_by_outer_id(outer_id)
        settings = cls._queryset.filter(channel=channel).first()
        return settings
