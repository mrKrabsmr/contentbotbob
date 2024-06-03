from django.contrib import admin

from apps.channels.models import Channel, ChannelSettings


class ChannelSettingsInline(admin.TabularInline):
    model = ChannelSettings
    extra = 1


@admin.register(Channel)
class ChannelAdmin(admin.ModelAdmin):
    inlines = [ChannelSettingsInline]
    list_display = ["owner", "name", "type", "resource", "status_on"]
