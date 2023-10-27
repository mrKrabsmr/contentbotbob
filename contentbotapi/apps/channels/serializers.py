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

    @atomic
    def create(self, validated_data):
        settings = validated_data.pop("settings")
        channel = Channel.objects.create(**validated_data)
        ChannelSettings.objects.create(channel=channel, **settings)

        return channel
