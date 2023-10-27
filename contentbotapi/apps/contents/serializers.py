from rest_framework import serializers
from rest_framework.serializers import ModelSerializer

from apps.contents.models import Content, ContentMedia


class NewsContentSerializer(serializers.Serializer):
    text = serializers.CharField()
    subject = serializers.CharField()
    source = serializers.CharField()
    img = serializers.CharField()


class ContentMediaSerializer(ModelSerializer):
    class Meta:
        model = ContentMedia
        fields = "__all__"


class ContentSerializer(ModelSerializer):
    images = ContentMediaSerializer(many=True, read_only=True)

    class Meta:
        model = Content
        fields = "__all__"
