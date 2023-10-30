from django.contrib import admin

from apps.channels.models import Channel, ChannelSettings, ChannelSettingsAllowedContentSource


class ChannelSettingsAllowedContentSourceInline(admin.TabularInline):
    model = ChannelSettingsAllowedContentSource
    extra = 1


@admin.register(ChannelSettings)
class ChannelSettingsAdmin(admin.ModelAdmin):
    inlines = (ChannelSettingsAllowedContentSourceInline,)


@admin.register(Channel)
class ChannelAdmin(admin.ModelAdmin):
    pass
