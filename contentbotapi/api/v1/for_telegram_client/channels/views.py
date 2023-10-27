import http

from rest_framework import mixins
from rest_framework.generics import GenericAPIView, UpdateAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from apps.channels.models import Channel, ChannelSettings
from apps.channels.serializers import ChannelSerializer, ChannelSettingsSerializer
from apps.channels.services import ChannelService as C_s
from apps.users.permissions import IsConfirmed, IsChannelOwner


class UserChannelListAPIView(APIView):
    permission_classes = [IsAuthenticated, IsConfirmed]

    def get(self, request, *args, **kwargs):
        channels = C_s.get_user_channels(request.user)
        serializer = ChannelSerializer(channels, many=True)

        return Response(serializer.data, http.HTTPStatus.OK)


class ChannelAPIView(mixins.UpdateModelMixin,
                     mixins.CreateModelMixin,
                     GenericAPIView):
    queryset = Channel.objects.all()
    serializer_class = ChannelSerializer
    permission_classes = [IsAuthenticated, IsConfirmed, IsChannelOwner]


class ChannelSettingsAPIView(UpdateAPIView):
    queryset = ChannelSettings.objects.all()
    serializer_class = ChannelSettingsSerializer
    permission_classes = [IsAuthenticated, IsConfirmed, IsChannelOwner]



