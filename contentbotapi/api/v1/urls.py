from django.urls import path
from rest_framework import routers
from rest_framework_simplejwt.views import TokenRefreshView

from api.v1.for_parser_client.contents.views import NewsContentAPIView
from api.v1.for_frontend_client.channels.views import ChannelViewSet, ChannelSettingsAPIView, \
    CheckPossibleAddChannelAPIView, ChannelListAPIView
from api.v1.for_frontend_client.contents.views import ContentAPIView
from api.v1.for_frontend_client.subscribes.views import SubscribeAPIView
from api.v1.for_frontend_client.users.views import RegisterAPIView, LoginAPIView, SendActivateCodeAPIView, \
    UserActivateAPIView, ProfileAPIView

channels_router = routers.SimpleRouter()
channels_router.register(r"channels", ChannelViewSet)

urlpatterns = [
    path("parser/contents/news/", NewsContentAPIView.as_view(), name="content_news"),

    path("register/", RegisterAPIView.as_view(), name="register"),
    path("login/", LoginAPIView.as_view(), name="login"),
    path("send/activate-code/", SendActivateCodeAPIView.as_view(), name="send_activate_code"),
    path("users/activate/", UserActivateAPIView.as_view(), name="user_activate"),
    path("refresh/", TokenRefreshView.as_view(), name="token_refresh"),
    path("profile/", ProfileAPIView.as_view(), name="profile"),
    path("channel-settings/<slug:outer_id>/", ChannelSettingsAPIView.as_view(), name="channel_settings"),
    path("list/channels/", ChannelListAPIView.as_view(), name="channel_list"),
    path("contents/for/<slug:outer_id>/", ContentAPIView.as_view(), name="content"),
    path("subscribes/", SubscribeAPIView.as_view(), name="subscribe"),
    path("check-add/channels/", CheckPossibleAddChannelAPIView.as_view(), name="check add channel")
] + channels_router.urls
