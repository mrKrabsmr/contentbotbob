from apps.services import ServiceCore
from apps.subscribes.models import Subscribe, UserSubscribe


class SubscribeService(ServiceCore):
    _queryset = Subscribe.objects.all()

    @classmethod
    def get_subscribe(cls, sub_id):
        return cls._queryset.filter(id=sub_id).first()


class UserSubscribeService(ServiceCore):
    _queryset = UserSubscribe.objects.all()

    @classmethod
    def check_user_already_have_subscibe(cls, user):
        return cls._queryset.filter(user=user).exists()

    @classmethod
    def activate_user_subscribe(cls, user, sub):
        cls._queryset.create(
            user=user,
            subscribe=sub
        )
