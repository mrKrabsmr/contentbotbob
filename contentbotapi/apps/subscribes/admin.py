from django.contrib import admin

from apps.subscribes.models import Subscribe, UserSubscribe


@admin.register(Subscribe)
class SubscribeAdmin(admin.ModelAdmin):
    list_display = ["name", "max_use_channels", "price_rub"]


@admin.register(UserSubscribe)
class UserSubscribeAdmin(admin.ModelAdmin):
    list_display = ["user", "subscribe", "bought_at"]

    def get_queryset(self, request):
        queryset = super().get_queryset(request)
        return queryset.select_related("user", "subscribe")
