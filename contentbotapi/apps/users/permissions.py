from rest_framework.permissions import BasePermission

from apps.channels.models import Channel, ChannelSettings


class IsConfirmed(BasePermission):
    """
    Сhecking the user's email confirmation
    """
    def has_permission(self, request, view):
        return request.user.is_confirmed


class IsChannelOwner(BasePermission):
    """
    Checking whether the user trying to influence the Channel or Channel Settings is the owner of the channel
    """
    def has_object_permission(self, request, view, obj):
        if isinstance(obj, Channel):
            return obj.owner == request.user

        if isinstance(obj, ChannelSettings):
            return obj.channel.owner == request.user

        return False
