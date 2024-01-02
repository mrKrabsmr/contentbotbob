from django.contrib import admin

from apps.channels.models import Channel, ChannelSettings


@admin.register(ChannelSettings)
class ChannelSettingsAdmin(admin.ModelAdmin):
    pass


@admin.register(Channel)
class ChannelAdmin(admin.ModelAdmin):
    pass
