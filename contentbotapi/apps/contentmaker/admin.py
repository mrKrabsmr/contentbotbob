from django.contrib import admin

from apps.contentmaker.models import Channel, ChannelSettingsAllowedContentSource, ChannelSettings, ChannelContent, \
    ChannelContentMedia


@admin.register(Channel)
class ChannelAdmin(admin.ModelAdmin):
    search_fields = ("name",)
    list_filter = ("type", "resource")
    list_display = ("name", "type", "resource")


class ChannelSettingsAllowedContentSourceInline(admin.TabularInline):
    model = ChannelSettingsAllowedContentSource
    extra = 0


@admin.register(ChannelSettings)
class ChannelSettingsAdmin(admin.ModelAdmin):
    list_select_related = ("channel",)
    inlines = (ChannelSettingsAllowedContentSourceInline,)


class ChannelContentMediaInline(admin.TabularInline):
    model = ChannelContentMedia
    extra = 0


@admin.register(ChannelContent)
class ChannelContentAdmin(admin.ModelAdmin):
    list_select_related = ("channel",)
    list_display = ("channel", "created_at")
    inlines = (ChannelContentMediaInline,)
