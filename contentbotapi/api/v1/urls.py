from django.urls import path
from rest_framework_simplejwt.views import TokenRefreshView

from api.v1.for_parser_client.contents.views import NewsContentAPIView
from api.v1.for_telegram_client.channels.views import UserChannelListAPIView, ChannelAPIView, ChannelSettingsAPIView
from api.v1.for_telegram_client.contents.views import ContentAPIView
from api.v1.for_telegram_client.subscribes.views import ActivateSubscribeAPIView
from api.v1.for_telegram_client.users.views import RegisterAPIView, LoginAPIView, SendActivateCodeAPIView, \
    UserActivateAPIView

urlpatterns = [
    path("parser/content/news/", NewsContentAPIView.as_view(), name="content_news"),

    path("register/", RegisterAPIView.as_view(), name="register"),
    path("login/", LoginAPIView.as_view(), name="login"),
    path("send/activate-code/", SendActivateCodeAPIView.as_view(), name="send_activate_code"),
    path("users/activate/", UserActivateAPIView.as_view(), name="user_activate"),
    path("refresh/", TokenRefreshView.as_view(), name="token_refresh"),

    path("channels/", UserChannelListAPIView.as_view(), name="channel_list"),
    path("channels/", ChannelAPIView.as_view()),
    path("channel-settings/", ChannelSettingsAPIView.as_view()),
    path("contents/", ContentAPIView.as_view(), name="content"),
    path("subscribes/", ActivateSubscribeAPIView.as_view(), name="activate_subscribe")
]

