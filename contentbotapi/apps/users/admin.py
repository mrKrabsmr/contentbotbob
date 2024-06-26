from django.contrib import admin
from django.contrib.auth.models import Group

from apps.users.models import User

admin.site.unregister(Group)


@admin.register(User)
class UserAdmin(admin.ModelAdmin):
    list_display = ["username", "is_confirmed"]
