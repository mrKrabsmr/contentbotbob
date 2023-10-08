from rest_framework import serializers


class NewsContentSerializer(serializers.Serializer):
    text = serializers.CharField()
    subject = serializers.CharField()
    source = serializers.CharField()
    img = serializers.CharField()

