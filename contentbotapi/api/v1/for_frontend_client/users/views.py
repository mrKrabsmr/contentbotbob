import http

from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.tokens import RefreshToken

from apps.subscribes.serializers import SubscribeSerializer, UserSubscribeSerializer
from apps.subscribes.services import UserSubscribeService as US_s
from apps.text import text_messages
from apps.users.models import User
from apps.users.serializers import RegisterSerializer, LoginSerializer
from apps.users.services import UserCodeActivationService as UCA_s


class RegisterAPIView(APIView):

    def post(self, request, *args, **kwargs):
        serializer = RegisterSerializer(data=request.data)

        if serializer.is_valid():
            serializer.save()
            UCA_s.send_activation_code(serializer.instance)

            return Response({"message": text_messages["sent_email"]}, http.HTTPStatus.CREATED)

        return Response(serializer.errors, http.HTTPStatus.BAD_REQUEST)


class LoginAPIView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = LoginSerializer(data=request.data)

        if serializer.is_valid():
            user = User.objects.filter(username=serializer.data["username"]).first()
            if user and user.check_password(serializer.data["password"]):
                token = RefreshToken.for_user(user)

                data = {
                    **serializer.data,
                    "access_token": str(token.access_token),
                    "refresh_token": str(token)
                }

                return Response(data, http.HTTPStatus.OK)
            return Response({"message": "пользователь с такими данными не найден"}, http.HTTPStatus.NOT_FOUND)
        return Response(serializer.errors, http.HTTPStatus.BAD_REQUEST)


class SendActivateCodeAPIView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request, *args, **kwargs):
        UCA_s.send_activation_code(request.user)
        return Response({"result": "ok"}, http.HTTPStatus.OK)


class UserActivateAPIView(APIView):
    permission_classes = [IsAuthenticated]

    def post(self, request, *args, **kwargs):
        if not UCA_s.check_user_have_code(request.user):
            return Response({"message": "activation code not found"}, http.HTTPStatus.NOT_FOUND)

        activation_code = UCA_s.get_object(request)
        if activation_code:
            activation_code.user.is_confirmed = True
            activation_code.user.save()
            activation_code.delete()

            return Response({"message": "success"}, http.HTTPStatus.OK)
        return Response({"message": "invalid code"}, http.HTTPStatus.BAD_REQUEST)


class ProfileAPIView(APIView):
    permission_classes = [IsAuthenticated]

    def get(self, request, *args, **kwargs):
        data = {
            "username": request.user.username,
            "email": request.user.email,
            "registered_at": request.user.date_joined.date(),
            "is_confirmed": request.user.is_confirmed,
            "subscribe": None
        }

        user_subscribe = US_s.get_user_subscribe(request.user)
        if user_subscribe:
            subscribe_serializer = SubscribeSerializer(user_subscribe.subscribe)

            sub_data = {
                "name": subscribe_serializer.data["name"],
                "max_use_channels": subscribe_serializer.data["max_use_channels"],
            }

            data["subscribe"] = sub_data

        return Response(data, http.HTTPStatus.OK)
