from django.urls import path
from rest_framework_simplejwt.views import TokenRefreshView

from api.for_server_client.contents.views import NewsContentAPIView

urlpatterns = [
    path('parser/content/news/', NewsContentAPIView.as_view(), name='content_news'),
    path('refresh/', TokenRefreshView.as_view(), name='token_refresh'),
]
