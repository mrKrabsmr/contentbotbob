from rest_framework.serializers import ModelSerializer

from apps.subscribes.models import Subscribe, UserSubscribe


class SubscribeSerializer(ModelSerializer):
    class Meta:
        model = Subscribe
        fields = "__all__"


class UserSubscribeSerializer(ModelSerializer):
    class Meta:
        model = UserSubscribe
        fields = "__all__"
