from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.tokens import RefreshToken

from apps.users import text_messages
from apps.users.models import UserCodeActivation, User
from apps.users.serializers import RegisterSerializer, LoginSerializer
from apps.users.utils import create_or_change_activation_code


class RegisterAPIView(APIView):

    def post(self, request, *args, **kwargs):
        serializer = RegisterSerializer(data=request.data)

        if serializer.is_valid():
            serializer.save()
            create_or_change_activation_code(serializer.instance)

            return Response({"message": text_messages["sent_email"]})

        return Response(serializer.error_messages, 400)


class LoginAPIView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = LoginSerializer(data=request.data)

        if serializer.is_valid():
            user = User.objects.get(username=serializer["username"])
            if user and user.check_password(serializer.data["password"]):
                token = RefreshToken.for_user(user)
                serializer.data["access_token"] = str(token.access_token)
                serializer.data["refresh_token"] = str(token)

                return Response(serializer.data, 200)
            return Response({"message": "invalid data"}, 403)
        return Response(serializer.errors, 400)


class SendActivateCodeAPIView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request, *args, **kwargs):
        create_or_change_activation_code(request.user)


class UserActivateAPIView(APIView):
    def post(self, request, *args, **kwargs):
        activation_code = UserCodeActivation.objects.filter(code=request.data["code"]).first()
        if activation_code:
            activation_code.user.is_confirmed = True
            activation_code.user.save()
            activation_code.delete()

            return Response({"message": "success"}, 200)
        return Response({"message": "activation code not found"}, 404)

