from django.contrib import admin

from apps.subscribes.models import Subscribe, UserSubscribe


@admin.register(Subscribe)
class SubscribeAdmin(admin.ModelAdmin):
    pass


@admin.register(UserSubscribe)
class UserSubscribeAdmin(admin.ModelAdmin):
    pass
