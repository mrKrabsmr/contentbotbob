from apps.users.models import UserCodeActivation
from apps.users.utils import get_activation_code, send_code


class UserCodeActivationService:
    _queryset = UserCodeActivation.objects.all()

    @classmethod
    def send_activation_code(cls, user):
        code_activation, created = cls._queryset.get_or_create(user=user)
        if not created:
            code_activation.code = get_activation_code()
            code_activation.save()

        send_code.delay(
            email=code_activation.user.email,
            code=code_activation.code
        )

    @classmethod
    def check_user_have_code(cls, user):
        return cls._queryset.filter(user=user).exists()

    @classmethod
    def get_object(cls, request):
        return cls._queryset.filter(user=request.user, code=request.data.get("code")).first()
