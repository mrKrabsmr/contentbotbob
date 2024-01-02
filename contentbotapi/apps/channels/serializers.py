from django.db.transaction import atomic
from rest_framework.serializers import ModelSerializer

from apps.channels.models import Channel, ChannelSettings


class ChannelSettingsSerializer(ModelSerializer):
    class Meta:
        model = ChannelSettings
        fields = "__all__"


class ChannelSerializer(ModelSerializer):
    settings = ChannelSettingsSerializer()

    class Meta:
        model = Channel
        fields = "__all__"

    def _get_user(self):
        return self.context.get("request").user

    @atomic
    def create(self, validated_data):
        settings = validated_data.pop("settings")
        user = self._get_user()
        channel = Channel.objects.create(**validated_data, owner=user)
        ChannelSettings.objects.create(channel=channel, **settings)

        return channel
