import http

from rest_framework import mixins
from rest_framework.generics import UpdateAPIView
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework.viewsets import GenericViewSet

from apps.channels.models import Channel, ChannelSettings
from apps.channels.serializers import ChannelSerializer, ChannelSettingsSerializer
from apps.channels.services import ChannelService as C_s
from apps.users.permissions import IsConfirmed, IsChannelOwner


class ChannelViewSet(mixins.UpdateModelMixin,
                     mixins.CreateModelMixin,
                     mixins.ListModelMixin,
                     GenericViewSet):
    queryset = Channel.objects.all()
    serializer_class = ChannelSerializer
    permission_classes = [IsAuthenticated, IsConfirmed, IsChannelOwner]

    def list(self, request, *args, **kwargs):
        channels = C_s.get_user_channels(request.user)
        serializer = ChannelSerializer(channels, many=True)

        return Response(serializer.data, http.HTTPStatus.OK)

    def create(self, request, *args, **kwargs):
        user = request.user
        if not C_s.check_user_may_add_channel(user):
            return Response({"message": "Нет прав на добавление канала"}, http.HTTPStatus.FORBIDDEN)

        return super().create(request, *args, **kwargs)


class ChannelSettingsAPIView(UpdateAPIView):
    queryset = ChannelSettings.objects.all()
    serializer_class = ChannelSettingsSerializer
    permission_classes = [IsAuthenticated, IsConfirmed, IsChannelOwner]


class CheckPossibleAddChannelAPIView(APIView):
    permission_classes = [IsAuthenticated, IsConfirmed]

    def get(self, request, *args, **kwargs):
        user = request.user
        may_add = C_s.check_user_may_add_channel(user)

        if may_add:
            return Response({"result": True}, http.HTTPStatus.OK)

        return Response({"result": False}, http.HTTPStatus.FORBIDDEN)
