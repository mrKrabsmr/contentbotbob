from rest_framework.serializers import ModelSerializer

from apps.subscribes.models import Subscribe


class SubstribeSerilaizer(ModelSerializer):
    class Meta:
        model = Subscribe
        fields = "__all__"
