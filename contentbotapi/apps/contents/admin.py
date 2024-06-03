from django.contrib import admin

from apps.contents.models import ContentMedia, Content


class ContentMediaInline(admin.TabularInline):
    model = ContentMedia
    extra = 1


@admin.register(Content)
class ContentAdmin(admin.ModelAdmin):
    inlines = (ContentMediaInline,)
    list_display = ["channel", "rating", "created_at"]

