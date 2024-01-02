from rest_framework import serializers
from rest_framework.serializers import ModelSerializer

from apps.contents.models import Content, ContentMedia


class NewsContentSerializer(serializers.Serializer):
    text = serializers.CharField()
    types = serializers.ListField(
        child=serializers.CharField()
    )
    img = serializers.ListField(
        child=serializers.CharField(required=False, allow_null=True, allow_blank=True),
        required=False,
        allow_null=True
    )


class ContentMediaSerializer(ModelSerializer):
    class Meta:
        model = ContentMedia
        exclude = ("content",)


class ContentSerializer(ModelSerializer):
    images = ContentMediaSerializer(many=True, read_only=True)

    class Meta:
        model = Content
        fields = ("id", "text", "images")
